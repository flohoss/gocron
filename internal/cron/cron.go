package cron

import (
	"github.com/go-co-op/gocron/v2"
)

func New() *Cron {
	s, _ := gocron.NewScheduler()
	return &Cron{scheduler: s}
}

type Cron struct {
	scheduler gocron.Scheduler
}

func (c *Cron) Add(cronString string, cmd func()) {
	c.scheduler.NewJob(gocron.CronJob(cronString, false), gocron.NewTask(cmd))
}

func (c *Cron) Run() {
	c.scheduler.Start()
}
