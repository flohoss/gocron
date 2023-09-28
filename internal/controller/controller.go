package controller

import (
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/env"
	"gitlab.unjx.de/flohoss/gobackup/internal/notify"
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
	service, err := database.MigrateDatabase(*notify.NewNotificationService(env.NtfyEndpoint, env.NtfyToken, env.NtfyTopic), env.Identifier)
	if err != nil {
		zap.L().Fatal("cannot connect to database", zap.Error(err))
	}
	database.SetupEventChannel()
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

type commands func(*database.Job, *database.Run)

func (c *Controller) runAllJobs(fn commands) {
	notify.SendHealthcheck(c.env.HealthcheckURL, c.env.HealthcheckUUID, "/start")
	jobs := c.service.GetJobs()
	for i := 0; i < len(jobs); i++ {
		c.runJob(fn, &jobs[i])
	}
	notify.SendHealthcheck(c.env.HealthcheckURL, c.env.HealthcheckUUID, "")
}

func (c *Controller) runJob(fn commands, job *database.Job) {
	run := database.Run{JobID: job.ID}
	c.service.CreateOrUpdate(&run)
	setupResticEnvVariables(job)
	fn(job, &run)
	removeResticEnvVariables(job)
	run.EndTime = time.Now().UnixMilli()
	c.service.CreateOrUpdate(&run)
}

func (c *Controller) runBackups() {
	c.runAllJobs(func(job *database.Job, run *database.Run) { c.runBackup(job, run) })
}

func (c *Controller) runPrunes() {
	c.runAllJobs(func(job *database.Job, run *database.Run) { c.runPrune(job, run) })
}

func (c *Controller) runChecks() {
	c.runAllJobs(func(job *database.Job, run *database.Run) { c.runCheck(job, run) })
}
