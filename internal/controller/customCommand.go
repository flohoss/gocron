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
	CustomCommand    string `json:"custom_command" validate:"omitempty"`
	LocalDirectory   string `json:"local_directory" validate:"omitempty,dir" example:"/"`
	ResticRemote     string `json:"restic_remote" validate:"omitempty" example:"rclone:pcloud:Backups/gitea"`
	PasswordFilePath string `json:"password_file_path" validate:"omitempty,file" example:"/secrets/.resticpwd"`
}

// @Schemes
// @Tags		commands
// @Accept		json
// @Param		command	body	CommandBody	true	"Command body"
// @Success	200
// @Failure	400	{object}	echo.HTTPError
// @Failure	404	{object}	echo.HTTPError
// @Router		/commands [post]
func (c *Controller) RunCommand(ctx echo.Context) error {
	cmdBody := new(CommandBody)
	if err := ctx.Bind(cmdBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(cmdBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch cmdBody.Command {
	case "restore":
		go c.restoreRepository(ctx, cmdBody)
		return ctx.NoContent(http.StatusOK)
	}

	job := c.service.GetJob(cmdBody.JobID)
	if job.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "job not found")
	}

	switch cmdBody.Command {
	case "run":
		go c.runJob(func(job *database.Job, run *database.Run) { c.runBackup(job, run) }, job, "backup")
	case "prune":
		go c.runJob(func(job *database.Job, run *database.Run) { c.runPrune(job, run) }, job, "pruning")
	case "check":
		go c.runJob(func(job *database.Job, run *database.Run) { c.runCheck(job, run) }, job, "check")
	case "custom":
		go c.runJob(func(job *database.Job, run *database.Run) {
			c.execute(ExecuteContext{runId: run.ID, logType: uint64(database.LogCustomCommand), errLogSeverity: uint64(database.LogError), successLog: true}, "restic", strings.Split(cmdBody.CustomCommand, " ")...)
		}, job, "custom command")
	}
	return ctx.NoContent(http.StatusOK)
}
