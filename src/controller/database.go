package controller

import (
	"gobackup/models"
	"os"

	"go.uber.org/zap"
)

func (c *Controller) databaseBackup(job *models.Job) {
	if job.DatabaseType == models.NoDatabase {
		return
	}
	var stderr, stdout []byte
	var err error
	password, _ := os.ReadFile(job.LocalDirectory + PasswordFile)
	switch job.DatabaseType {
	case models.PostgreSQL:
		stderr, stdout, err = c.executeDockerCmdSeperateOut("exec", "-e", "PGPASSWORD="+string(password), job.DatabaseContainer, "pg_dump", job.DatabaseName, "--username="+job.DatabaseUser)
	case models.MariaDB:
		stderr, stdout, err = c.executeDockerCmdSeperateOut("exec", job.DatabaseContainer, "mysqldump", "--user="+job.DatabaseUser, "--password="+string(password), "--lock-tables", "--databases", job.DatabaseName)
	}
	if err != nil {
		c.addLogEntry(models.Log{JobID: job.ID, Type: models.Warn, Topic: models.Database, Message: string(stderr) + string(stdout)}, job.Description)
		zap.L().Warn("cannot backup database", zap.ByteString("stderr", stderr), zap.ByteString("stdout", stdout))
		return
	}
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Database, Message: "Database dump successful"}, job.Description)
	file, err := os.OpenFile(job.LocalDirectory+"/.dbBackup.sql", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		zap.L().Error("cannot open .dbBackup.sql", zap.Error(err))
	}
	defer file.Close()
	_, err = file.Write(stdout)
	if err != nil {
		zap.L().Error("cannot write database", zap.Error(err))
	}
	zap.L().Debug("Database dump successful", zap.String("file", file.Name()))
}
