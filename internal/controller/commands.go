package controller

import (
	"fmt"
	"os/exec"
	"strings"

	"gitlab.unjx.de/flohoss/gobackup/database"
)

type ExecuteContext struct {
	runId           uint64
	localDirectory  string
	logType         uint64
	errLogSeverity  uint64
	errMsgOverwrite string
	successLog      bool
}

func (c *Controller) execute(ctx ExecuteContext, program string, commands ...string) error {
	cmd := exec.Command(program, commands...)
	if ctx.localDirectory != "" {
		cmd.Dir = ctx.localDirectory
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.errMsgOverwrite != "" {
			c.service.CreateOrUpdate(&database.Log{
				RunID:         ctx.runId,
				LogTypeID:     ctx.logType,
				LogSeverityID: ctx.errLogSeverity,
				Message:       ctx.errMsgOverwrite,
			})
		} else {
			c.service.CreateOrUpdate(&database.Log{
				RunID:         ctx.runId,
				LogTypeID:     ctx.logType,
				LogSeverityID: ctx.errLogSeverity,
				Message:       string(out),
			})
		}
		return fmt.Errorf("%s", out)
	}
	if ctx.successLog {
		msg := string(out)
		if msg == "" {
			msg = program
			for _, str := range commands {
				msg += " " + str
			}
		}
		c.service.CreateOrUpdate(&database.Log{
			RunID:         ctx.runId,
			LogTypeID:     ctx.logType,
			LogSeverityID: uint64(database.LogInfo),
			Message:       msg,
		})
	}
	return nil
}

func (c *Controller) handlePreAndPostCommands(localDirectory string, cmds []database.Command, runId uint64) error {
	for _, cmd := range cmds {
		split := strings.Split(cmd.Command, " ")
		if len(split) >= 2 {
			err := c.execute(ExecuteContext{
				runId:          runId,
				localDirectory: localDirectory,
				logType:        uint64(database.LogCustomCommand),
				errLogSeverity: uint64(database.LogError),
				successLog:     true,
			}, split[0], split[1:]...)
			if err != nil {
				return err
			}
		} else {
			c.service.CreateOrUpdate(&database.Log{
				RunID:         runId,
				LogTypeID:     uint64(database.LogCustomCommand),
				LogSeverityID: uint64(database.LogWarning),
				Message:       fmt.Sprintf("command '%s' is missing parameters", cmd.Command),
			})
		}
	}
	return nil
}

func (c *Controller) runBackup(job *database.Job, run *database.Run) error {
	if !c.resticRepositoryExists(job, run) {
		if err := c.initResticRepository(job, run); err != nil {
			return err
		}
	}
	if err := c.handlePreAndPostCommands(job.LocalDirectory, job.PreCommands, run.ID); err != nil {
		return err
	}
	if err := c.execute(ExecuteContext{
		runId:          run.ID,
		logType:        uint64(database.LogRestic),
		errLogSeverity: uint64(database.LogError),
		successLog:     true,
	}, "restic", "backup", job.LocalDirectory, "--no-scan", "--compression", job.CompressionType.Compression); err != nil {
		return err
	}
	if err := c.handlePreAndPostCommands(job.LocalDirectory, job.PostCommands, run.ID); err != nil {
		return err
	}
	return nil
}

func (c *Controller) runPrune(job *database.Job, run *database.Run) error {
	if c.resticRepositoryExists(job, run) {
		if job.RetentionPolicy.ID == 1 {
			c.service.CreateOrUpdate(&database.Log{
				RunID:         run.ID,
				LogTypeID:     uint64(database.LogPrune),
				LogSeverityID: uint64(database.LogInfo),
				Message:       "keeping all snapshots, nothing to do...",
			})
			return nil
		}
		retPolicy := strings.Split(job.RetentionPolicy.Policy, " ")
		combined := append([]string{"forget", "--prune"}, retPolicy...)
		if err := c.execute(ExecuteContext{
			runId:          run.ID,
			logType:        uint64(database.LogPrune),
			errLogSeverity: uint64(database.LogError),
			successLog:     true,
		}, "restic", combined...); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runCheck(job *database.Job, run *database.Run) error {
	if c.resticRepositoryExists(job, run) {
		if job.RoutineCheck == 0 {
			c.service.CreateOrUpdate(&database.Log{
				RunID:         run.ID,
				LogTypeID:     uint64(database.LogCheck),
				LogSeverityID: uint64(database.LogInfo),
				Message:       "routine check is disabled",
			})
			return nil
		}
		if err := c.execute(ExecuteContext{
			runId:          run.ID,
			logType:        uint64(database.LogCheck),
			errLogSeverity: uint64(database.LogError),
			successLog:     true,
		}, "restic", "check", fmt.Sprintf("--read-data-subset=%d%%", job.RoutineCheck)); err != nil {
			return err
		}
	}
	return nil
}
