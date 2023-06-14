package controller

import (
	"fmt"
	"gobackup/models"
	"os"

	"go.uber.org/zap"
)

func setupResticEnvVariables(job *models.Job) {
	os.Setenv("RESTIC_REPOSITORY", getResticRepo(job))
	os.Setenv("RESTIC_PASSWORD_FILE", job.Remote.PasswordFile)
}

func removeResticEnvVariables(job *models.Job) {
	os.Unsetenv("RESTIC_REPOSITORY")
	os.Unsetenv("RESTIC_PASSWORD_FILE")
}

func getResticRepo(job *models.Job) string {
	return fmt.Sprintf("%s/%s", job.Remote.Repository, job.RemoteDirectory)
}

func (c *Controller) resticRepositoryExists(job *models.Job) bool {
	_, err := c.executeCmd("restic", "snapshots", "-q")
	if err != nil {
		msg := "no existing repository found"
		c.addLogEntry(models.Log{JobID: job.ID, Type: models.Warn, Topic: models.Restic, Message: msg}, job.Description)
		zap.L().Warn(msg)
		return false
	}
	return true
}

func (c *Controller) initResticRepository(job *models.Job) error {
	out, err := c.executeCmd("restic", "init")
	if err != nil {
		c.addLogEntry(models.Log{JobID: job.ID, Type: models.Error, Topic: models.Restic, Message: string(out)}, job.Description)
		zap.L().Error("cannot initialize repository", zap.String("job", job.Description), zap.ByteString("msg", out))
		return fmt.Errorf("%s", out)
	}
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: string(out)}, job.Description)
	zap.L().Debug("repository initialized", zap.String("job", job.Description), zap.ByteString("msg", out))
	return nil
}

func (c *Controller) runResticCommand(job *models.Job, cmd ...string) error {
	return c.runCommand(job, "restic", cmd...)
}
