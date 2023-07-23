package controller

import (
	"fmt"
	"os"
	"os/exec"

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
	_, err := exec.Command("restic", "snapshots", "-q").Output()
	if err != nil {
		c.createLog(&database.Log{
			RunID:         run.ID,
			LogTypeID:     uint64(database.LogTypeBackup),
			LogSeverityID: uint64(database.LogSeverityWarning),
			Message:       "no existing repository found",
		})
		return false
	}
	return true
}

func (c *Controller) initResticRepository(job *database.Job, run *database.Run) error {
	out, err := exec.Command("restic", "init").CombinedOutput()
	if err != nil {
		c.createLog(&database.Log{
			RunID:         run.ID,
			LogTypeID:     uint64(database.LogTypeBackup),
			LogSeverityID: uint64(database.LogSeverityError),
			Message:       string(out),
		})
		return fmt.Errorf("%s", out)
	}
	c.createLog(&database.Log{
		RunID:         run.ID,
		LogTypeID:     uint64(database.LogTypeBackup),
		LogSeverityID: uint64(database.LogSeverityInfo),
		Message:       string(out),
	})
	return nil
}
