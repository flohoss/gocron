package controller

import (
	"fmt"
	"gobackup/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type JobsData struct {
	Title      string
	Jobs       []models.Job
	JobRunning bool
}

type LogJobData struct {
	Title string
}

type FormJobData struct {
	Title   string
	Job     *models.Job
	Options []models.SelectOption
	Remotes []models.SelectOption
	Errors  map[string][]string
}

func (c *Controller) RenderJobs(ctx echo.Context) error {
	var jobs []models.Job
	c.orm.Preload("Remote").Order("Description").Find(&jobs)
	return ctx.Render(http.StatusOK, "jobs", JobsData{Title: c.env.Identifier + " - Jobs", Jobs: jobs, JobRunning: c.JobRunning})
}

func (c *Controller) CreateJobConfiguration(ctx echo.Context) error {
	job := new(models.Job)
	err := c.bindToRequest(ctx, job)
	if err == nil {
		data := FormJobData{Title: c.env.Identifier + " - New Job", Options: models.DatabaseTypes(), Remotes: c.RemoteOptions(), Job: job, Errors: make(map[string][]string)}
		if job.ID != 0 {
			data.Title = c.env.Identifier + " - Edit Job"
			c.setFormPassword(job)
			if job.DatabaseType == models.NoDatabase {
				job.DatabaseContainer = ""
				job.DatabaseName = ""
				job.DatabaseUser = ""
			}
		}

		if errors := c.createOrUpdate(ctx, job); len(errors) != 0 {
			data.Errors = errors
			return ctx.Render(http.StatusOK, "jobsForm", data)
		}

		c.savePassword(job, ctx.FormValue("database_password"))
		if strings.Contains(data.Title, "New") {
			c.SSE.CreateStream(fmt.Sprintf("%s%d", EventLog, job.ID))
			c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Backup, Message: "Job created"}, job.Description)
		}
	}
	return c.RenderJobs(ctx)
}

func (c *Controller) RenderJobsLog(ctx echo.Context) error {
	var logs []models.Log
	result := c.orm.Where("job_id = ?", ctx.QueryParam("id")).Limit(1).Find(&logs)
	if result.RowsAffected == 0 {
		ctx.Redirect(http.StatusTemporaryRedirect, "/jobs")
	}
	return ctx.Render(http.StatusOK, "jobsLog", LogJobData{Title: c.env.Identifier + " - Logs"})
}

func (c *Controller) RenderJobForm(ctx echo.Context) error {
	job := new(models.Job)
	id := ctx.QueryParam("id")
	title := "New Job"
	if id != "" {
		c.orm.Limit(1).Find(&job, id)
		if job.ID != 0 {
			title = "Edit Job"
			c.setFormPassword(job)
		}
	}
	return ctx.Render(http.StatusOK, "jobsForm", FormJobData{Title: c.env.Identifier + " - " + title, Job: job, Options: models.DatabaseTypes(), Remotes: c.RemoteOptions()})
}

func (c *Controller) DeleteJobConfiguration(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := c.orm.Delete(&models.Job{}, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.SSE.RemoveStream("log" + id)
	return ctx.NoContent(http.StatusOK)
}

func (c *Controller) StartBackup(ctx echo.Context) error {
	if c.JobRunning {
		return ctx.NoContent(http.StatusServiceUnavailable)
	}
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}
	job := new(models.Job)
	c.orm.Limit(1).Preload("Remote").Find(job, id)
	if job.ID == 0 {
		return ctx.NoContent(http.StatusNotFound)
	}
	go c.runJob(func(job *models.Job) {
		if err := c.runBackup(job); err != nil {
			c.updateJobStatus(job, models.Stopped)
		}
	}, job)
	return ctx.NoContent(http.StatusOK)
}

func (c *Controller) getJobOptions() []models.SelectOption {
	jobs := []models.Job{}
	options := []models.SelectOption{}
	c.orm.Order("description").Find(&jobs)
	for _, job := range jobs {
		options = append(options, models.SelectOption{Name: job.Description, Value: int(job.ID)})
	}
	return options
}
