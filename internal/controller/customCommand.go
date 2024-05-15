package controller

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/database"
)

type CommandBody struct {
	Command          string `json:"command" validate:"required"`
	JobID            uint64 `json:"job_id" validate:"omitempty,number"`
	CustomCommand    string `json:"custom_command" validate:"omitempty,ascii"`
	LocalDirectory   string `json:"local_directory" validate:"omitempty,dir" example:"/"`
	ResticRemote     string `json:"restic_remote" validate:"omitempty" example:"rclone:pcloud:Backups/gitea"`
	PasswordFilePath string `json:"password_file_path" validate:"omitempty,file" example:"/secrets/.resticpwd"`
}

//	@Schemes
//	@Tags		commands
//	@Accept		json
//	@Param		command	body	CommandBody	true	"Command body"
//	@Success	200
//	@Failure	400	{object}	echo.HTTPError
//	@Failure	404	{object}	echo.HTTPError
//	@Router		/commands [post]
func (c *Controller) RunCommand(ctx echo.Context) error {
	jobs := c.service.GetJobs()
	for _, j := range jobs {
		if j.Status == database.LogNone {
			return echo.NewHTTPError(http.StatusBadRequest, "another job already running")
		}
	}
	cmdBody := new(CommandBody)
	if err := ctx.Bind(cmdBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	jsonBlob, err := c.service.ValidateRequestBinding(ctx, cmdBody)
	if err != nil {
		return ctx.JSONBlob(http.StatusBadRequest, jsonBlob)
	}

	switch cmdBody.Command {
	case "restore":
		go c.restoreRepository(cmdBody)
		return ctx.NoContent(http.StatusOK)
	}

	job := c.service.GetJob(cmdBody.JobID)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}

	switch cmdBody.Command {
	case "run":
		go c.runJob(func(job *database.Job, run *database.Run) { c.runBackup(job, run) }, job)
	case "prune":
		go c.runJob(func(job *database.Job, run *database.Run) { c.runPrune(job, run) }, job)
	case "check":
		go c.runJob(func(job *database.Job, run *database.Run) { c.runCheck(job, run) }, job)
	case "custom":
		split := strings.Split(cmdBody.CustomCommand, " ")
		if len(split) >= 2 {
			go c.runJob(func(job *database.Job, run *database.Run) {
				c.execute(ExecuteContext{
					runId:          run.ID,
					localDirectory: job.LocalDirectory,
					logType:        database.LogCustom,
					errLogSeverity: database.LogError,
					successLog:     true,
				}, split[0], split[1:]...)
			}, job)
		} else {
			return ctx.NoContent(http.StatusBadRequest)
		}
	}
	return ctx.NoContent(http.StatusOK)
}
