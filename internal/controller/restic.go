package controller

import (
	"os"

	"github.com/labstack/echo/v4"
	"gitlab.unjx.de/flohoss/gobackup/database"
)

func setupResticEnvVariables(job *database.Job) {
	os.Setenv("RESTIC_REPOSITORY", job.ResticRemote)
	os.Setenv("RESTIC_PASSWORD_FILE", job.PasswordFilePath)
}

func removeResticEnvVariables(job *database.Job) {
	os.Unsetenv("RESTIC_REPOSITORY")
	os.Unsetenv("RESTIC_PASSWORD_FILE")
}

func (c *Controller) resticRepositoryExists(job *database.Job, run *database.Run) bool {
	err := c.execute(ExecuteContext{
		runId:           run.ID,
		logType:         database.LogRestic,
		errLogSeverity:  database.LogWarning,
		errMsgOverwrite: "no existing repository found",
	}, "restic", "snapshots", "-q")
	return err == nil
}

func (c *Controller) initResticRepository(job *database.Job, run *database.Run) error {
	return c.execute(ExecuteContext{
		runId:          run.ID,
		logType:        database.LogRestic,
		errLogSeverity: database.LogError,
		successLog:     true,
	}, "restic", "init")
}

func (c *Controller) restoreRepository(ctx echo.Context, cmdBody *CommandBody) {
	if cmdBody.ResticRemote == "" {
		c.service.CreateOrUpdate(&database.SystemLog{
			LogSeverity: database.LogError,
			Message:     "no restic remote provided",
		})
	}
	if cmdBody.PasswordFilePath == "" {
		c.service.CreateOrUpdate(&database.SystemLog{
			LogSeverity: database.LogError,
			Message:     "no password file provided",
		})
	}
	if cmdBody.LocalDirectory == "" {
		cmdBody.LocalDirectory = "/"
	}
	c.executeSystem(ExecuteContext{
		errLogSeverity: database.LogError,
		successLog:     true,
	}, "restic", "-r", cmdBody.ResticRemote, "restore", "latest", "--target", cmdBody.LocalDirectory, "--password-file", cmdBody.PasswordFilePath)
}
