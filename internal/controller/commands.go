package controller

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gitlab.unjx.de/flohoss/gobackup/database"
	"gitlab.unjx.de/flohoss/gobackup/internal/notify"
)

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
	c.createLog(&database.Log{
		RunID:         run.ID,
		LogTypeID:     uint64(database.LogTypeBackup),
		LogSeverityID: uint64(database.LogSeverityInfo),
		Message:       "run started",
	})
	setupResticEnvVariables(job)
	fn(job, &run)
	removeResticEnvVariables(job)
	run.EndTime = time.Now().UnixMilli()
	c.createLog(&database.Log{
		RunID:         run.ID,
		LogTypeID:     uint64(database.LogTypeBackup),
		LogSeverityID: uint64(database.LogSeverityInfo),
		Message:       "run stopped",
	})
	c.service.CreateOrUpdate(&run)
}

func (c *Controller) runBackups() {
	c.runAllJobs(func(job *database.Job, run *database.Run) { c.runBackup(job, run) })
}

func (c *Controller) runPrunes() {
	c.runAllJobs(func(job *database.Job, run *database.Run) { c.runPrune(job, run) })
}

func (c *Controller) runChecks() {
	c.runAllJobs(func(job *database.Job, run *database.Run) { c.runCheck(job, run, c.env.DefaultSubset, false) })
}

func (c *Controller) runBackup(job *database.Job, run *database.Run) error {
	if !c.resticRepositoryExists(job, run) {
		if err := c.initResticRepository(job, run); err != nil {
			return err
		}
	}
	if err := c.handleCommands(job.PreCommands, run.ID); err != nil {
		return err
	}
	if err := c.execute(ExecuteContext{
		runId:           run.ID,
		logType:         uint64(database.LogTypeBackup),
		errLogSeverity:  uint64(database.LogSeverityError),
		errMsgOverwrite: "",
		successLog:      true,
	}, "restic", "backup", job.LocalDirectory, "--no-scan", "--compression", job.CompressionType.Compression); err != nil {
		return err
	}
	if err := c.handleCommands(job.PostCommands, run.ID); err != nil {
		return err
	}
	return nil
}

func (c *Controller) runPrune(job *database.Job, run *database.Run) error {
	if c.resticRepositoryExists(job, run) {
		retPolicy := strings.Split(job.RetentionPolicy.Policy, " ")
		combined := append([]string{"forget"}, retPolicy...)
		combined = append(combined, []string{"--prune"}...)
		if out, err := exec.Command("restic", combined...).CombinedOutput(); err != nil {
			return fmt.Errorf("%s", out)
		}
	}
	return nil
}

func (c *Controller) runCheck(job *database.Job, run *database.Run, subset uint, overwrite bool) error {
	if c.resticRepositoryExists(job, run) {
		if out, err := exec.Command("restic", "check", fmt.Sprintf("--read-data-subset=%d%%", subset)).CombinedOutput(); err != nil {
			return fmt.Errorf("%s", out)
		}
	}
	return nil
}
