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

	"github.com/enescakir/emoji"
	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

const (
	DATE_FORMAT = "2006-01-02 15:04:05"
)

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

func formatTime(startTime int64) string {
	startSeconds := startTime / 1000
	t := time.Unix(startSeconds, 0).Local()
	return t.Format(DATE_FORMAT)
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
		s.Add(sTime, func() {
			js.ExecuteJobs(jobQueues[sTime])
		})
	}

	if s.DeleteRunsAfterDays > 0 {
		s.Add("0 0 * * *", func() {
			queries.DeleteRuns(ctx, time.Now().AddDate(0, 0, -int(s.DeleteRunsAfterDays)).UnixMilli())
		})
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
	for i := range jobs {
		js.ExecuteJob(&jobs[i])
		names = append(names, jobs[i].Name)
	}
	js.Notify.Send(fmt.Sprintf("Backup finished %v", emoji.PartyPopper), fmt.Sprintf("Time: %s\nJobs: %s", time.Now().Format(time.RFC1123), strings.Join(names, ", ")), log.INFO)
}

func (js *JobService) ExecuteJob(job *jobs.Job) {
	ctx := context.Background()
	dbJob, _ := js.ListJob(job.ID, 4)

	runView := js.startRun(ctx, dbJob)

	envs, _ := js.Queries.ListEnvsByJobID(ctx, job.ID)
	keys := []string{}
	for _, e := range envs {
		os.Setenv(e.Key, commands.ExtractVariable(e.Value))
		keys = append(keys, e.Key)
	}
	js.writeLog(ctx, dbJob, runView.ID, Debug, fmt.Sprintf("Setting environment variables: \"%s\"", strings.Join(keys, ", ")))

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
		js.writeLog(ctx, dbJob, runView.ID, Debug, msg)
		out, err := commands.ExecuteCommand(program, args, command.FileOutput)
		severity = Info
		if err != nil {
			severity = Error
			js.Notify.Send(fmt.Sprintf("Error - %s", job.Name), fmt.Sprintf("Command: \"%s\"\nResult: \"%s\"", cmd, out), log.ERROR)
		}
		if out == "" {
			out = "Done - No output"
		}
		js.writeLog(ctx, dbJob, runView.ID, severity, out)
		if err != nil {
			runView.StatusID = Stopped.Int64()
			break
		}
	}

	for _, e := range envs {
		os.Unsetenv(e.Key)
	}

	js.endRun(ctx, dbJob, runView)
}

func (js *JobService) ListJobs() []jobs.JobsView {
	resultSet, _ := js.Queries.GetJobsView(context.Background())
	jobsAmount := len(resultSet)
	for i := range jobsAmount {
		resultSet[i].Runs, _ = js.Queries.GetRunsView(context.Background(), jobs.GetRunsViewParams{JobID: resultSet[i].ID, Limit: 3})
	}
	return resultSet
}

func (js *JobService) ListJob(id string, limit int64) (*jobs.JobsView, error) {
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

	jobView.Runs, _ = js.Queries.GetRunsView(context.Background(), jobs.GetRunsViewParams{JobID: job.ID, Limit: limit})
	amount := len(jobView.Runs)
	for i := 0; i < amount; i++ {
		logs, _ := js.Queries.ListLogsByRunID(context.Background(), jobView.Runs[i].ID)
		jobView.Runs[i].Logs = logs
	}

	return &jobView, err
}

func (js *JobService) refreshLogs(jobView *jobs.JobsView, run *jobs.Run, newLog *jobs.Log) {
	amount := len(jobView.Runs)
	// most likely the last run
	for i := amount - 1; i >= 0; i-- {
		if jobView.Runs[i].ID == run.ID {
			createdAtSeconds := newLog.CreatedAt / 1000
			t := time.Unix(createdAtSeconds, 0).Local()
			formattedTime := t.Format(DATE_FORMAT)

			jobView.Runs[i].Logs = append(jobView.Runs[i].Logs, jobs.ListLogsByRunIDRow{
				CreatedAt:     newLog.CreatedAt,
				RunID:         newLog.RunID,
				SeverityID:    newLog.SeverityID,
				Message:       newLog.Message,
				CreatedAtTime: formattedTime,
			})

			break
		}
	}
}

func (js *JobService) startRun(ctx context.Context, dbJob *jobs.JobsView) *jobs.RunsView {
	run, _ := js.Queries.CreateRun(ctx, jobs.CreateRunParams{
		JobID:     dbJob.ID,
		StatusID:  int64(Running),
		StartTime: time.Now().UnixMilli(),
	})

	runView := &jobs.RunsView{
		ID:           run.ID,
		JobID:        dbJob.ID,
		StatusID:     run.StatusID,
		StartTime:    run.StartTime,
		EndTime:      run.EndTime,
		FmtStartTime: formatTime(run.StartTime),
		Logs:         nil,
	}
	dbJob.Runs = append(dbJob.Runs, *runView)

	js.Events.SendEvent(true, dbJob)
	// prepare run to be finished if no error is set
	runView.StatusID = Finished.Int64()
	return runView
}

func (js *JobService) endRun(ctx context.Context, dbJob *jobs.JobsView, runView *jobs.RunsView) {
	run, _ := js.Queries.UpdateRun(ctx, jobs.UpdateRunParams{
		StatusID: runView.StatusID,
		EndTime:  sql.NullInt64{Int64: time.Now().UnixMilli(), Valid: true},
		ID:       runView.ID,
	})

	amount := len(dbJob.Runs)
	// most likely the last run
	for i := amount - 1; i >= 0; i-- {
		if dbJob.Runs[i].ID == run.ID {
			dbJob.Runs[i].FmtEndTime.String = formatTime(run.EndTime.Int64)
			dbJob.Runs[i].FmtEndTime.Valid = true
			dbJob.Runs[i].Duration.Int64 = run.EndTime.Int64 - run.StartTime
			dbJob.Runs[i].Duration.Valid = true
			break
		}
	}

	js.Events.SendEvent(true, dbJob)
}

func (js *JobService) writeLog(ctx context.Context, dbJob *jobs.JobsView, runId int64, severity Severity, message string) {
	newLog, _ := js.Queries.CreateLog(ctx, jobs.CreateLogParams{
		CreatedAt:  time.Now().UnixMilli(),
		RunID:      runId,
		SeverityID: int64(severity),
		Message:    message,
	})

	js.refreshLogs(dbJob, &jobs.Run{ID: runId}, &newLog)
	js.Events.SendEvent(false, dbJob)
}
