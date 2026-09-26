package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
)

func New() *Scheduler {
	s := cron.New()
	s.Start()
	return &Scheduler{scheduler: s}
}

type Scheduler struct {
	scheduler *cron.Cron
}

func (c *Scheduler) Stop() context.Context {
	return c.scheduler.Stop()
}

func (c *Scheduler) Add(cronString string, cmd func()) error {
	_, err := c.scheduler.AddFunc(cronString, cmd)
	return err
}
