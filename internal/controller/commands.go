package controller

import (
	"fmt"
	"os/exec"
	"strings"

	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/notify"
)

type commands func(*database.Job)

func (c *Controller) runAllJobs(fn commands) {
	notify.SendHealthcheck(c.env.HealthcheckURL, c.env.HealthcheckUUID, "/start")
	jobs := c.service.GetJobs()
	for i := 0; i < len(jobs); i++ {
		c.runJob(func(job *database.Job) {
			fn(&jobs[i])
		}, &jobs[i])
	}
	notify.SendHealthcheck(c.env.HealthcheckURL, c.env.HealthcheckUUID, "")
}

func (c *Controller) runJob(fn commands, job *database.Job) {
	setupResticEnvVariables(job)
	fn(job)
	removeResticEnvVariables(job)
}

func (c *Controller) runBackups() {
	c.runAllJobs(func(job *database.Job) {
		c.runBackup(job)
	})
}

func (c *Controller) runPrunes() {
	c.runAllJobs(func(job *database.Job) {
		c.runPrune(job)
	})
}

func (c *Controller) runChecks() {
	c.runAllJobs(func(job *database.Job) {
		c.runCheck(job, c.env.DefaultSubset, false)
	})
}

func (c *Controller) runBackup(job *database.Job) error {
	if !c.resticRepositoryExists(job) {
		if err := c.initResticRepository(job); err != nil {
			return err
		}
	}
	if out, err := exec.Command("restic", "backup", job.LocalDirectory, "--no-scan", "--compression", job.CompressionType.Compression).CombinedOutput(); err != nil {
		return fmt.Errorf("%s", out)
	}
	return nil
}

func (c *Controller) runPrune(job *database.Job) error {
	if c.resticRepositoryExists(job) {
		retPolicy := strings.Split(job.RetentionPolicy.Policy, " ")
		combined := append([]string{"forget"}, retPolicy...)
		combined = append(combined, []string{"--prune"}...)
		if out, err := exec.Command("restic", combined...).CombinedOutput(); err != nil {
			return fmt.Errorf("%s", out)
		}
	}
	return nil
}

func (c *Controller) runCheck(job *database.Job, subset uint, overwrite bool) error {
	if c.resticRepositoryExists(job) {
		if out, err := exec.Command("restic", "check", fmt.Sprintf("--read-data-subset=%d%%", subset)).CombinedOutput(); err != nil {
			return fmt.Errorf("%s", out)
		}
	}
	return nil
}
