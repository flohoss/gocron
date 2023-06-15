package controller

import (
	"os"

	"gitlab.unjx.de/flohoss/gobackup/internal/models"
)

const (
	PasswordFile        = "/.dbpass"
	PasswordPlaceHolder = "**********"
)

func (c *Controller) savePassword(job *models.Job, password string) {
	pwFile := job.LocalDirectory + PasswordFile
	if job.DatabaseType == models.NoDatabase {
		os.Remove(pwFile)
		return
	} else if password != PasswordPlaceHolder {
		file, _ := os.OpenFile(pwFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		defer file.Close()
		file.Write([]byte(password))
	}
}

func (c *Controller) setFormPassword(job *models.Job) {
	if _, err := os.Stat(job.LocalDirectory + PasswordFile); err == nil {
		job.DatabasePassword = PasswordPlaceHolder
	}
}
