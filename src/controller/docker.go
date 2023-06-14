package controller

import (
	"fmt"
	"gobackup/models"
	"path/filepath"

	"go.uber.org/zap"
)

func (c *Controller) stopDocker(job *models.Job) {
	if job.DockerRestart {
		c.runDockerCommand(job, "stop")
	}
}

func (c *Controller) startDocker(job *models.Job) {
	if job.DockerRestart {
		c.runDockerCommand(job, "up", "-d", "--no-recreate")
	}
}

func checkIfDockerComposeExists(job *models.Job) (file string, err error) {
	files, err := filepath.Glob(job.LocalDirectory + "/docker-compose.*")
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("docker compose file does not exist")
	}
	return files[0], nil
}

func (c *Controller) runDockerCommand(job *models.Job, cmd ...string) error {
	file, err := checkIfDockerComposeExists(job)
	if err != nil {
		c.addLogEntry(models.Log{JobID: job.ID, Type: models.Warn, Topic: models.Docker, Message: err.Error()}, job.Description)
		zap.L().Warn("cannot run docker compose command", zap.Error(err))
		return err
	}
	return c.runCommand(job, "docker", append([]string{"compose", "-f", file}, cmd...)...)
}
