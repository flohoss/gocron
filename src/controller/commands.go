package controller

import (
	"bytes"
	"fmt"
	"gobackup/models"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func getCommandTopic(cmdType string) models.LogTopic {
	switch cmdType {
	case "restic":
		return models.Restic
	case "docker":
		return models.Docker
	default:
		return models.Backup
	}
}

func (c *Controller) runCommand(job *models.Job, cmdType string, cmd ...string) error {
	out, err := c.executeCmd(cmdType, cmd...)
	if err != nil {
		c.addLogEntry(models.Log{JobID: job.ID, Type: models.Error, Topic: getCommandTopic(cmdType), Message: string(out)}, job.Description)
		zap.L().Error("cannot execute command", zap.ByteString("msg", out))
		return fmt.Errorf("%s", out)
	}
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: getCommandTopic(cmdType), Message: string(out)}, job.Description)
	zap.L().Debug("command executed", zap.ByteString("msg", out))
	return nil
}

func (c *Controller) runSystemCommand(cmdType string, cmd ...string) error {
	out, err := c.executeCmd(cmdType, cmd...)
	if err != nil {
		c.addSystemLogEntry(models.SystemLog{Type: models.Error, Topic: getCommandTopic(cmdType), Message: string(out)})
		zap.L().Error("cannot execute system command", zap.ByteString("msg", out))
		return fmt.Errorf("%s", out)
	}
	c.addSystemLogEntry(models.SystemLog{Type: models.Info, Topic: getCommandTopic(cmdType), Message: string(out)})
	zap.L().Debug("system command executed", zap.ByteString("msg", out))
	return nil
}

func (c *Controller) executeCmd(name string, commands ...string) ([]byte, error) {
	zap.L().Debug("executing command", zap.Strings("cmd", append([]string{name}, commands...)))
	var out []byte
	var err error
	// https://semgrep.dev/docs/cheat-sheets/go-command-injection/
	switch name {
	case "restic":
		out, err = exec.Command("restic", commands...).CombinedOutput()
	case "rclone":
		out, err = exec.Command("rclone", commands...).CombinedOutput()
	case "docker":
		out, err = exec.Command("docker", commands...).CombinedOutput()
	}
	return out, err
}

func (c *Controller) executeDockerCmdSeperateOut(commands ...string) ([]byte, []byte, error) {
	zap.L().Debug("executing command", zap.Strings("cmd", append([]string{"docker"}, commands...)))
	cmd := exec.Command("docker", commands...)
	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stderr.Bytes(), stdout.Bytes(), err
}

type commands func(*models.Job)

func (c *Controller) runAllJobs(fn commands) {
	c.setJobsRunning(true)
	c.sendHealthcheck("/start")
	var jobs []models.Job
	c.orm.Preload("Remote").Order("Description").Find(&jobs)
	for i := 0; i < len(jobs); i++ {
		setupResticEnvVariables(&jobs[i])
		c.updateJobStatus(&jobs[i], models.Running)
		fn(&jobs[i])
		c.updateJobStatus(&jobs[i], models.Success)
		removeResticEnvVariables(&jobs[i])
	}
	c.sendHealthcheck("")
	c.setJobsRunning(false)
}

func (c *Controller) runJob(fn commands, job *models.Job) {
	c.setJobsRunning(true)
	setupResticEnvVariables(job)
	c.updateJobStatus(job, models.Running)
	fn(job)
	c.updateJobStatus(job, models.Success)
	removeResticEnvVariables(job)
	c.setJobsRunning(false)
}

func (c *Controller) runBackups() {
	c.runAllJobs(func(job *models.Job) {
		if err := c.runBackup(job); err != nil {
			c.updateJobStatus(job, models.Stopped)
			c.startDocker(job)
		}
	})
}

func (c *Controller) runPrunes() {
	c.runAllJobs(func(job *models.Job) {
		if err := c.runPrune(job); err != nil {
			c.updateJobStatus(job, models.Stopped)
		}
	})
}

func (c *Controller) runChecks() {
	c.runAllJobs(func(job *models.Job) {
		if err := c.runCheck(job, c.env.DefaultSubset, false); err != nil {
			c.updateJobStatus(job, models.Stopped)
		}
	})
}

func (c *Controller) runBackup(job *models.Job) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Backup, Message: "starting backup"}, job.Description)
	c.databaseBackup(job)
	c.stopDocker(job)
	if !c.resticRepositoryExists(job) {
		if err := c.initResticRepository(job); err != nil {
			return err
		}
	}
	if err := c.runResticCommand(job, "backup", job.LocalDirectory, "--no-scan", "--compression", models.ResticCompressionType(job.Remote.CompressionType)); err != nil {
		return err
	}
	c.startDocker(job)
	return nil
}

func (c *Controller) runPrune(job *models.Job) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: "starting prune"}, job.Description)
	if c.resticRepositoryExists(job) {
		if job.Remote.RetentionPolicy == "" {
			msg := "no retention policy specified"
			c.addLogEntry(models.Log{JobID: job.ID, Type: models.Warn, Topic: models.Restic, Message: msg}, job.Description)
			zap.L().Warn(msg)
			return nil
		}
		retPolicy := strings.Split(job.Remote.RetentionPolicy, " ")
		combined := append([]string{"forget"}, retPolicy...)
		combined = append(combined, []string{"--prune"}...)
		if err := c.runResticCommand(job, combined...); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runCheck(job *models.Job, subset uint, overwrite bool) error {
	if !job.CheckResticRepo && !overwrite {
		return nil
	}
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: fmt.Sprintf("starting check with %d%% subset", subset)}, job.Description)
	if c.resticRepositoryExists(job) {
		if err := c.runResticCommand(job, "check", fmt.Sprintf("--read-data-subset=%d%%", subset)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runListSnapshots(job *models.Job) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: "listing snapshots"}, job.Description)
	if err := c.runResticCommand(job, "snapshots"); err != nil {
		return err
	}
	return nil
}

func (c *Controller) runRepairIndex(job *models.Job) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: "repairing index"}, job.Description)
	if c.resticRepositoryExists(job) {
		if err := c.runResticCommand(job, "repair", "index"); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runRepairSnapshots(job *models.Job, forget bool) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: "repairing snapshots"}, job.Description)
	if c.resticRepositoryExists(job) {
		cmd := []string{"repair", "snapshots"}
		if forget {
			cmd = append(cmd, "--forget")
		}
		if err := c.runResticCommand(job, cmd...); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runUnlock(job *models.Job, removeAll bool) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Restic, Message: "unlocking repo"}, job.Description)
	if c.resticRepositoryExists(job) {
		cmd := []string{"unlock"}
		if removeAll {
			cmd = append(cmd, "--remove-all")
		}
		if err := c.runResticCommand(job, cmd...); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) runLogs(job *models.Job, amount uint) error {
	c.addLogEntry(models.Log{JobID: job.ID, Type: models.Info, Topic: models.Docker, Message: "getting logs"}, job.Description)
	if err := c.runDockerCommand(job, "logs", "--tail", strconv.FormatUint(uint64(amount), 10)); err != nil {
		return err
	}
	return nil
}
