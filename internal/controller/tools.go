package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/internal/models"
)

type ToolsData struct {
	Title      string
	JobOptions []models.SelectOption
}

type DockerRequestBody struct {
	Command      string `json:"command" validate:"required"`
	JobID        uint   `json:"job_id" validate:"omitempty,number"`
	Amount       uint   `json:"amount" validate:"required,number"`
	AllImages    bool   `json:"all_images" validate:"omitempty"`
	PruneVolumes bool   `json:"prune_volumes" validate:"omitempty"`
}

type ResticRequestBody struct {
	Command          string `json:"command" validate:"required"`
	JobID            uint   `json:"job_id" validate:"omitempty,number"`
	Amount           uint   `json:"amount" validate:"required,number"`
	RemoteRepository string `json:"remote_repository" validate:"omitempty"`
	PasswordFile     string `json:"password_file" validate:"omitempty"`
	RepairForget     bool   `json:"repair_forget" validate:"omitempty"`
	UnlockRemoveAll  bool   `json:"unlock_remove_all" validate:"omitempty"`
}

func (c *Controller) RenderTools(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "tools", ToolsData{Title: c.env.Identifier + " - Tools", JobOptions: c.getJobOptions()})
}

func (c *Controller) DockerRequest(ctx echo.Context) error {
	params := new(DockerRequestBody)
	if err := ctx.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch params.Command {
	case "prune":
		cmd := []string{"system", "prune", "--force"}
		if params.AllImages {
			cmd = append(cmd, "--all")
		}
		if params.PruneVolumes {
			cmd = append(cmd, "--volumes")
		}
		go c.runSystemCommand("docker", cmd...)
		return ctx.NoContent(http.StatusOK)
	case "logs":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			return ctx.NoContent(http.StatusNotFound)
		}
		go c.runJob(func(job *models.Job) {
			if err := c.runLogs(job, params.Amount); err != nil {
				c.updateJobStatus(job, models.Stopped)
			}
		}, job)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.NoContent(http.StatusBadRequest)
}

func (c *Controller) ResticRequest(ctx echo.Context) error {
	if c.JobRunning {
		return ctx.NoContent(http.StatusServiceUnavailable)
	}
	params := new(ResticRequestBody)
	if err := ctx.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch params.Command {
	case "backup":
		go c.runBackups()
		return ctx.NoContent(http.StatusOK)
	case "forget":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			go c.runPrunes()
		} else {
			go c.runJob(func(job *models.Job) {
				if err := c.runPrune(job); err != nil {
					c.updateJobStatus(job, models.Stopped)
				}
			}, job)
		}
		return ctx.NoContent(http.StatusOK)
	case "check":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			c.runAllJobs(func(job *models.Job) {
				if err := c.runCheck(job, params.Amount, true); err != nil {
					c.updateJobStatus(job, models.Stopped)
				}
			})
		} else {
			go c.runJob(func(job *models.Job) {
				if err := c.runCheck(job, params.Amount, true); err != nil {
					c.updateJobStatus(job, models.Stopped)
				}
			}, job)
		}
		return ctx.NoContent(http.StatusOK)
	case "snapshots":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			return ctx.NoContent(http.StatusNotFound)
		}
		go c.runJob(func(job *models.Job) {
			if err := c.runListSnapshots(job); err != nil {
				c.updateJobStatus(job, models.Stopped)
			}
		}, job)
		return ctx.NoContent(http.StatusOK)
	case "repair-index":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			return ctx.NoContent(http.StatusNotFound)
		}
		go c.runJob(func(job *models.Job) {
			if err := c.runRepairIndex(job); err != nil {
				c.updateJobStatus(job, models.Stopped)
			}
		}, job)
		return ctx.NoContent(http.StatusOK)
	case "repair-snapshots":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			return ctx.NoContent(http.StatusNotFound)
		}
		go c.runJob(func(job *models.Job) {
			if err := c.runRepairSnapshots(job, params.RepairForget); err != nil {
				c.updateJobStatus(job, models.Stopped)
			}
		}, job)
		return ctx.NoContent(http.StatusOK)
	case "unlock":
		job := new(models.Job)
		c.orm.Limit(1).Preload("Remote").Find(job, params.JobID)
		if job.ID == 0 {
			return ctx.NoContent(http.StatusNotFound)
		}
		go c.runJob(func(job *models.Job) {
			if err := c.runUnlock(job, params.UnlockRemoveAll); err != nil {
				c.updateJobStatus(job, models.Stopped)
			}
		}, job)
		return ctx.NoContent(http.StatusOK)
	case "restore":
		if params.RemoteRepository == "" || params.PasswordFile == "" {
			return ctx.NoContent(http.StatusBadRequest)
		}
		c.addSystemLogEntry(models.SystemLog{Type: models.Info, Topic: models.Restic, Message: fmt.Sprintf("restoring repo '%s'", params.RemoteRepository)})
		go c.runSystemCommand("restic", "-r", params.RemoteRepository, "restore", "latest", "--target", "/", "--password-file", params.PasswordFile)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.NoContent(http.StatusBadRequest)
}
