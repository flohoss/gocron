package controller

import (
	"fmt"
	"os"
	"os/exec"

	"gitlab.unjx.de/flohoss/gobackup/database"
	"go.uber.org/zap"
)

func setupResticEnvVariables(job *database.Job) {
	os.Setenv("RESTIC_REPOSITORY", job.ResticRemote)
	os.Setenv("RESTIC_PASSWORD_FILE", job.PasswordFilePath)
}

func removeResticEnvVariables(job *database.Job) {
	os.Unsetenv("RESTIC_REPOSITORY")
	os.Unsetenv("RESTIC_PASSWORD_FILE")
}

func (c *Controller) resticRepositoryExists(job *database.Job) bool {
	_, err := exec.Command("restic", "snapshots", "-q").Output()
	if err != nil {
		msg := "no existing repository found"
		zap.L().Warn(msg)
		return false
	}
	return true
}

func (c *Controller) initResticRepository(job *database.Job) error {
	out, err := exec.Command("restic", "init").CombinedOutput()
	if err != nil {
		zap.L().Error("cannot initialize repository", zap.String("job", job.Description), zap.ByteString("msg", out))
		return fmt.Errorf("%s", out)
	}
	zap.L().Debug("repository initialized", zap.String("job", job.Description), zap.ByteString("msg", out))
	return nil
}
