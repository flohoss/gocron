package services

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"github.com/robfig/cron/v3"

	"gitlab.unjx.de/flohoss/gocron/config"
	"gitlab.unjx.de/flohoss/gocron/internal/commands"
	"gitlab.unjx.de/flohoss/gocron/internal/events"
	"gitlab.unjx.de/flohoss/gocron/internal/notify"
	"gitlab.unjx.de/flohoss/gocron/internal/scheduler"
	"gitlab.unjx.de/flohoss/gocron/services/jobs"
)

//go:embed jobs.sql
var ddl string

func generateID(input string) string {
	var result strings.Builder

	// Iterate over each character in the input string
	for _, ch := range input {
		// Convert to lowercase and check if the character is alphanumeric or a space
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			result.WriteRune(unicode.ToLower(ch)) // Convert to lowercase and add it
		} else if ch == ' ' {
			result.WriteRune('_') // Replace space with underscore
		}
	}

	return result.String()
}

func NewJobService(dbName string, config *config.Config, s *scheduler.Scheduler, notify *notify.Notifier) (*JobService, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", dbName+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := jobs.New(db)
	initEnums(queries, ctx)

	err = createUpdateOrDeleteJob(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	err = createUpdateOrDeleteEnvs(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	err = createUpdateOrDeleteCommands(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	// no need for config any longer as all information is in db
	config = nil

	var jobNames = []string{}
	var jobQueues = make(map[string][]jobs.Job)
	js := &JobService{Queries: queries, Notify: notify, Scheduler: s}

	dbJobs, _ := queries.ListJobs(ctx)

	for _, job := range dbJobs {
		jobNames = append(jobNames, job.ID)
		jobQueues[job.Cron] = append(jobQueues[job.Cron], job)
	}

	for sTime := range jobQueues {
		s.Add(sTime, func(sTime string) func() {
			return func() {
				js.ExecuteJobs(jobQueues[sTime])
			}
		}(sTime))
	}

	js.Events = events.New(jobNames, func(streamID string, sub *sse.Subscriber) {
		js.Events.SendEvent(js.IsIdle(), nil)
	})

	return js, nil
}

type JobService struct {
	Queries   *jobs.Queries
	Notify    *notify.Notifier
	Scheduler *scheduler.Scheduler
	Events    *events.Event
}

func initEnums(queries *jobs.Queries, ctx context.Context) {
	severities, _ := queries.ListSeverities(ctx)
	if len(severities) == 0 {
		queries.CreateSeverity(ctx, Debug.String())
		queries.CreateSeverity(ctx, Info.String())
		queries.CreateSeverity(ctx, Warning.String())
		queries.CreateSeverity(ctx, Error.String())
	}
	status, _ := queries.ListStatus(ctx)
	if len(status) == 0 {
		queries.CreateStatus(ctx, Running.String())
		queries.CreateStatus(ctx, Stopped.String())
		queries.CreateStatus(ctx, Finished.String())
	}
}

func createUpdateOrDeleteJob(ctx context.Context, queries *jobs.Queries, config *config.Config) error {
	dbJobs, _ := queries.ListJobs(ctx)

	existingJobs := make(map[string]bool)
	for _, j := range dbJobs {
		existingJobs[j.ID] = true
	}

	for _, j := range config.Jobs {
		jobID := generateID(j.Name)
		if _, exists := existingJobs[jobID]; exists {
			queries.UpdateJob(ctx, jobs.UpdateJobParams{
				ID:   jobID,
				Name: j.Name,
				Cron: j.Cron,
			})
		} else {
			queries.CreateJob(ctx, jobs.CreateJobParams{
				ID:   jobID,
				Name: j.Name,
				Cron: j.Cron,
			})
		}
		delete(existingJobs, jobID)
	}

	for id := range existingJobs {
		queries.DeleteJob(ctx, id)
	}

	return nil
}

func createUpdateOrDeleteEnvs(ctx context.Context, queries *jobs.Queries, config *config.Config) error {
	queries.DeleteEnvs(ctx)
	var created int64 = 0
	for _, job := range config.Jobs {
		for _, env := range job.Envs {
			created++
			_, err := queries.CreateEnv(ctx, jobs.CreateEnvParams{
				ID:    created,
				JobID: generateID(job.Name),
				Key:   env.Key,
				Value: env.Value,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createUpdateOrDeleteCommands(ctx context.Context, queries *jobs.Queries, config *config.Config) error {
	queries.DeleteCommands(ctx)
	var created int64 = 0
	for _, job := range config.Jobs {
		for _, command := range job.Commands {
			created++
			create := jobs.CreateCommandParams{
				ID:      created,
				JobID:   generateID(job.Name),
				Command: command.Command,
			}
			if command.FileOutput != "" {
				create.FileOutput = sql.NullString{
					String: command.FileOutput,
					Valid:  true,
				}
			}
			_, err := queries.CreateCommand(ctx, create)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (js *JobService) GetQueries() *jobs.Queries {
	return js.Queries
}

func (js *JobService) GetParser() *cron.Parser {
	return js.Scheduler.GetParser()
}

func (js *JobService) GetHandler() echo.HandlerFunc {
	return js.Events.GetHandler()
}

func (js *JobService) IsIdle() bool {
	res, _ := js.Queries.IsIdle(context.Background())
	return res == 1
}

func (js *JobService) ExecuteJobs(jobs []jobs.Job) {
	if len(jobs) == 0 {
		jobs, _ = js.Queries.ListJobs(context.Background())
	}
	names := []string{}
	for i := 0; i < len(jobs); i++ {
		js.ExecuteJob(&jobs[i])
		names = append(names, jobs[i].Name)
	}
	if js.Notify.SendMessageOnSuccess {
		js.Notify.Send("Backup finished", fmt.Sprintf("Time: %s\nJobs: %s", time.Now().Format(time.RFC1123), strings.Join(names, ", ")), notify.DEFAULT, []string{"tada"})
	}
}

func (js *JobService) ExecuteJob(job *jobs.Job) {
	ctx := context.Background()
	dbJob, _ := js.ListJob(job.ID)

	run := js.startRun(ctx, dbJob)
	status := Finished

	envs, _ := js.Queries.ListEnvsByJobID(ctx, job.ID)
	keys := []string{}
	for _, e := range envs {
		os.Setenv(e.Key, commands.ExtractVariable(e.Value))
		keys = append(keys, e.Key)
	}
	js.writeLog(ctx, dbJob, run.ID, Debug, fmt.Sprintf("Setting environment variables: \"%s\"", strings.Join(keys, ", ")))

	cmds, _ := js.Queries.ListCommandsByJobID(ctx, job.ID)
	for _, command := range cmds {
		severity := Debug
		program, args := commands.PrepareCommand(command.Command)
		cmd := program
		if len(args) != 0 {
			cmd += " " + strings.Join(args, " ")
		}
		msg := fmt.Sprintf("Executing command: \"%s\"", cmd)
		if command.FileOutput.Valid {
			msg = fmt.Sprintf("Executing command (output to file): \"%s\"", cmd)
		}
		js.writeLog(ctx, dbJob, run.ID, Debug, msg)
		out, err := commands.ExecuteCommand(program, args, command.FileOutput)
		severity = Info
		if err != nil {
			severity = Error
			js.Notify.Send(fmt.Sprintf("Error - %s", job.Name), fmt.Sprintf("Command: \"%s\"\nResult: \"%s\"", cmd, out), notify.URGENT, []string{"rotating_light"})
		}
		if out == "" {
			out = "Done - No output"
		}
		js.writeLog(ctx, dbJob, run.ID, severity, out)
		if err != nil {
			status = Stopped
			break
		}
	}

	for _, e := range envs {
		os.Unsetenv(e.Key)
	}

	js.endRun(ctx, dbJob, run.ID, status)
}

func (js *JobService) ListJobs() []jobs.JobsView {
	resultSet, _ := js.Queries.GetJobsView(context.Background())
	jobsAmount := len(resultSet)
	for i := 0; i < jobsAmount; i++ {
		resultSet[i].Runs, _ = js.Queries.GetRunsView(context.Background(), resultSet[i].ID)
	}
	return resultSet
}

func (js *JobService) ListJob(id string) (*jobs.JobsView, error) {
	job, err := js.Queries.GetJob(context.Background(), id)
	if err != nil {
		return nil, err
	}

	jobView := jobs.JobsView{
		ID:   job.ID,
		Name: job.Name,
		Cron: job.Cron,
		Runs: nil,
	}

	jobView.Runs, _ = js.Queries.GetRunsView(context.Background(), job.ID)
	amount := len(jobView.Runs)
	for i := 0; i < amount; i++ {
		logs, _ := js.Queries.ListLogsByRunID(context.Background(), jobView.Runs[i].ID)
		jobView.Runs[i].Logs = logs
	}

	return &jobView, err
}

func (js *JobService) refreshRunsAndLogs(jobView *jobs.JobsView) {
	jobView.Runs, _ = js.Queries.GetRunsView(context.Background(), jobView.ID)
	js.refreshLogs(jobView)
}

func (js *JobService) refreshLogs(jobView *jobs.JobsView) {
	amount := len(jobView.Runs)
	for i := 0; i < amount; i++ {
		logs, _ := js.Queries.ListLogsByRunID(context.Background(), jobView.Runs[i].ID)
		jobView.Runs[i].Logs = logs
	}
}

func (js *JobService) startRun(ctx context.Context, dbJob *jobs.JobsView) *jobs.Run {
	run, _ := js.Queries.CreateRun(ctx, jobs.CreateRunParams{
		JobID:     dbJob.ID,
		StatusID:  int64(Running),
		StartTime: time.Now().UnixMilli(),
	})
	js.refreshRunsAndLogs(dbJob)
	js.Events.SendEvent(true, dbJob)
	return &run
}

func (js *JobService) endRun(ctx context.Context, dbJob *jobs.JobsView, runId int64, status Status) {
	js.Queries.UpdateRun(ctx, jobs.UpdateRunParams{
		StatusID: int64(status),
		EndTime:  sql.NullInt64{Int64: time.Now().UnixMilli(), Valid: true},
		ID:       runId,
	})
	js.refreshRunsAndLogs(dbJob)
	js.Events.SendEvent(true, dbJob)
}

func (js *JobService) writeLog(ctx context.Context, dbJob *jobs.JobsView, runId int64, severity Severity, message string) {
	js.Queries.CreateLog(ctx, jobs.CreateLogParams{
		CreatedAt:  time.Now().UnixMilli(),
		RunID:      runId,
		SeverityID: int64(severity),
		Message:    message,
	})
	js.refreshLogs(dbJob)
	js.Events.SendEvent(false, dbJob)
}
