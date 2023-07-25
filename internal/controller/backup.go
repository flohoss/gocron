package controller

import (
	"fmt"
	"os/exec"
	"strings"

	"gitlab.unjx.de/flohoss/gobackup/database"
)

type ExecuteContext struct {
	runId           uint64
	logType         uint64
	errLogSeverity  uint64
	errMsgOverwrite string
	successLog      bool
}

func (c *Controller) execute(ctx ExecuteContext, program string, commands ...string) error {
	out, err := exec.Command(program, commands...).CombinedOutput()
	if err != nil {
		if ctx.errMsgOverwrite != "" {
			c.createLog(&database.Log{
				RunID:         ctx.runId,
				LogTypeID:     ctx.logType,
				LogSeverityID: ctx.errLogSeverity,
				Message:       ctx.errMsgOverwrite,
			})
		} else {
			c.createLog(&database.Log{
				RunID:         ctx.runId,
				LogTypeID:     ctx.logType,
				LogSeverityID: ctx.errLogSeverity,
				Message:       string(out),
			})
		}
		return fmt.Errorf("%s", out)
	}
	if ctx.successLog {
		c.createLog(&database.Log{
			RunID:         ctx.runId,
			LogTypeID:     ctx.logType,
			LogSeverityID: uint64(database.LogSeverityInfo),
			Message:       string(out),
		})
	}
	return nil
}

func (c *Controller) handleCommands(cmds []database.Command, runId uint64) error {
	for _, cmd := range cmds {
		split := strings.Split(cmd.Command, " ")
		if len(split) >= 2 {
			err := c.execute(ExecuteContext{
				runId:           runId,
				logType:         uint64(database.LogTypeBackup),
				errLogSeverity:  uint64(database.LogSeverityError),
				errMsgOverwrite: "",
				successLog:      true,
			}, split[0], split[1:]...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
