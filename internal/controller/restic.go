package controller

import (
	"os"

	"gitlab.unjx.de/flohoss/gobackup/database"
)

func setupResticEnvVariables(job *database.Job) {
	os.Setenv("RESTIC_REPOSITORY", job.ResticRemote)
	os.Setenv("RESTIC_PASSWORD_FILE", job.PasswordFilePath)
}

func removeResticEnvVariables() {
	os.Unsetenv("RESTIC_REPOSITORY")
	os.Unsetenv("RESTIC_PASSWORD_FILE")
}

func (c *Controller) resticRepositoryExists(run *database.Run) bool {
	return c.execute(ExecuteContext{
		runId:           run.ID,
		logType:         database.LogRestic,
		errLogSeverity:  database.LogWarning,
		errMsgOverwrite: "no existing repository found",
		successLog:      true,
	}, "restic", "cat", "config") != nil
}

func (c *Controller) initResticRepository(run *database.Run) error {
	return c.execute(ExecuteContext{
		runId:          run.ID,
		logType:        database.LogRestic,
		errLogSeverity: database.LogError,
		successLog:     true,
	}, "restic", "init")
}

func (c *Controller) restoreRepository(cmdBody *CommandBody) {
	if cmdBody.LocalDirectory == "" {
		cmdBody.LocalDirectory = "/"
	}
	cmd := []string{"-r", cmdBody.ResticRemote, "restore", "latest", "--target", cmdBody.LocalDirectory, "--password-file", cmdBody.PasswordFilePath}
	stringCmd := "restic"
	for _, c := range cmd {
		stringCmd += " " + c
	}
	c.service.CreateOrUpdate(&database.SystemLog{
		LogSeverity: database.LogInfo,
		Message:     "cmd: " + stringCmd,
	})
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
	c.executeSystem(ExecuteContext{
		errLogSeverity: database.LogError,
		successLog:     true,
	}, "restic", cmd...)
}
