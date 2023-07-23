package controller

import (
	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"go.uber.org/zap"
)

type Controller struct {
	service  *database.Service
	env      *env.Config
	schedule *cron.Cron
}

type IndexData struct {
	Title string
	Jobs  []database.Job
}

func NewController(env *env.Config) *Controller {
	service, err := database.MigrateDatabase()
	if err != nil {
		zap.L().Fatal("cannot connect to database", zap.Error(err))
	}
	ctrl := Controller{service: service, env: env}
	ctrl.setupSchedule()

	return &ctrl
}

func (c *Controller) setupSchedule() {
	c.schedule = cron.New()
	if c.env.BackupCron != "" {
		c.schedule.AddFunc(c.env.BackupCron, func() {
			c.runBackups()
		})
	}
	if c.env.CleanupCron != "" {
		c.schedule.AddFunc(c.env.CleanupCron, func() {
			c.runPrunes()
		})
	}
	if c.env.CheckCron != "" {
		c.schedule.AddFunc(c.env.CheckCron, func() {
			c.runChecks()
		})
	}
	c.schedule.Start()
}
