package controller

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gitlab.unjx.de/flohoss/gobackup/database"
)

type ExecuteContext struct {
	runId           uint64
	localDirectory  string
	fileOutput      string
	logType         database.LogType
	errLogSeverity  database.LogSeverity
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
				RunID:       ctx.runId,
				LogType:     ctx.logType,
				LogSeverity: ctx.errLogSeverity,
				Message:     ctx.errMsgOverwrite,
			})
		} else {
			c.service.CreateOrUpdate(&database.Log{
				RunID:       ctx.runId,
				LogType:     ctx.logType,
				LogSeverity: ctx.errLogSeverity,
				Message:     string(out),
			})
		}
		return fmt.Errorf("%s", out)
	}

	if ctx.fileOutput != "" {
		file, err := os.OpenFile(ctx.localDirectory+"/"+ctx.fileOutput, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err == nil {
			defer file.Close()
			file.Write(out)
		} else {
			c.service.CreateOrUpdate(&database.Log{
				RunID:       ctx.runId,
				LogType:     ctx.logType,
				LogSeverity: ctx.errLogSeverity,
				Message:     "cannot write command to file",
			})
		}
	}

	if ctx.successLog {
		msg := string(out)
		if ctx.fileOutput != "" || msg == "" {
			msg = program
			for _, str := range commands {
				msg += " " + str
			}
			msg = database.AnonymisePasswords(msg)
		}
		c.service.CreateOrUpdate(&database.Log{
			RunID:       ctx.runId,
			LogType:     ctx.logType,
			LogSeverity: database.LogInfo,
			Message:     msg,
		})
	}
	return nil
}

func (c *Controller) executeSystem(ctx ExecuteContext, program string, commands ...string) error {
	cmd := exec.Command(program, commands...)
	if ctx.localDirectory != "" {
		cmd.Dir = ctx.localDirectory
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.errMsgOverwrite != "" {
			c.service.CreateOrUpdate(&database.SystemLog{
				LogSeverity: ctx.errLogSeverity,
				Message:     ctx.errMsgOverwrite,
			})
		} else {
			c.service.CreateOrUpdate(&database.SystemLog{
				LogSeverity: ctx.errLogSeverity,
				Message:     string(out),
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
		c.service.CreateOrUpdate(&database.SystemLog{
			LogSeverity: database.LogInfo,
			Message:     msg,
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
				fileOutput:     cmd.FileOutput,
				logType:        database.LogCustom,
				errLogSeverity: database.LogError,
				successLog:     true,
			}, split[0], split[1:]...)
			if err != nil {
				return err
			}
		} else {
			c.service.CreateOrUpdate(&database.Log{
				RunID:       runId,
				LogType:     database.LogCustom,
				LogSeverity: database.LogWarning,
				Message:     fmt.Sprintf("command '%s' is missing parameters", cmd.Command),
			})
		}
	}
	return nil
}

func (c *Controller) runBackup(job *database.Job, run *database.Run) error {
	if !c.resticRepositoryExists(run) {
		if err := c.initResticRepository(run); err != nil {
			return err
		}
	}
	if err := c.handlePreAndPostCommands(job.LocalDirectory, job.PreCommands, run.ID); err != nil {
		return err
	}
	c.execute(ExecuteContext{
		runId:          run.ID,
		logType:        database.LogRestic,
		errLogSeverity: database.LogError,
		successLog:     true,
	}, "restic", "backup", job.LocalDirectory, "--no-scan", "--compression", database.CompressionTypeInfoMap[job.CompressionType].Command)
	return c.handlePreAndPostCommands(job.LocalDirectory, job.PostCommands, run.ID)
}

func (c *Controller) runPrune(job *database.Job, run *database.Run) error {
	if c.resticRepositoryExists(run) {
		if job.RetentionPolicy == database.KeepAll {
			c.service.CreateOrUpdate(&database.Log{
				RunID:       run.ID,
				LogType:     database.LogPrune,
				LogSeverity: database.LogInfo,
				Message:     "keeping all snapshots, nothing to do...",
			})
			return nil
		}
		retPolicy := strings.Split(database.RetentionPolicyInfoMap[job.RetentionPolicy].Command, " ")
		combined := append([]string{"forget", "--prune"}, retPolicy...)
		if err := c.execute(ExecuteContext{
			runId:          run.ID,
			logType:        database.LogPrune,
			errLogSeverity: database.LogError,
			successLog:     true,
		}, "restic", combined...); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runCheck(job *database.Job, run *database.Run) error {
	if c.resticRepositoryExists(run) {
		if job.RoutineCheck == 0 {
			c.service.CreateOrUpdate(&database.Log{
				RunID:       run.ID,
				LogType:     database.LogCheck,
				LogSeverity: database.LogInfo,
				Message:     "routine check is disabled",
			})
			return nil
		}
		if err := c.execute(ExecuteContext{
			runId:          run.ID,
			logType:        database.LogCheck,
			errLogSeverity: database.LogError,
			successLog:     true,
		}, "restic", "check", fmt.Sprintf("--read-data-subset=%d%%", job.RoutineCheck)); err != nil {
			return err
		}
	}
	return nil
}
