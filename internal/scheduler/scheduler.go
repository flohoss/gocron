package scheduler

import (
	"github.com/robfig/cron/v3"
)

func New(deleteRunsAfterDays int) *Scheduler {
	s := cron.New()
	s.Start()
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	return &Scheduler{scheduler: s, parser: parser, DeleteRunsAfterDays: deleteRunsAfterDays}
}

type Scheduler struct {
	scheduler           *cron.Cron
	parser              cron.Parser
	DeleteRunsAfterDays int
}

func (c *Scheduler) Add(cronString string, cmd func()) {
	c.scheduler.AddFunc(cronString, cmd)
}

func (c *Scheduler) GetParser() *cron.Parser { return &c.parser }
