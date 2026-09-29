package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

func New(location *time.Location) *Scheduler {
	s := cron.New(cron.WithLocation(location))
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

func (c *Scheduler) Location() *time.Location {
	return c.scheduler.Location()
}
