package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"github.com/robfig/cron/v3"

	"gitlab.unjx.de/flohoss/gocron/config"
	"gitlab.unjx.de/flohoss/gocron/internal/commands"
	"gitlab.unjx.de/flohoss/gocron/internal/events"
	"gitlab.unjx.de/flohoss/gocron/internal/scheduler"
	"gitlab.unjx.de/flohoss/gocron/services/jobs"
)

const (
	DATE_FORMAT = "2006-01-02 15:04:05"
)

func formatTime(startTime int64) string {
	startSeconds := startTime / 1000
	t := time.Unix(startSeconds, 0).Local()
	return t.Format(DATE_FORMAT)
}

func NewJobService() (*JobService, error) {
	s := scheduler.New()
	queries, err := setupSQLite()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var cronJobs = make(map[string]map[string]config.Job)
	js := &JobService{Queries: queries, Scheduler: s}

	for name, job := range config.GetJobs() {
		cron := job.Cron
		if _, exists := cronJobs[cron]; !exists {
			cronJobs[cron] = make(map[string]config.Job)
		}
		cronJobs[cron][name] = job
	}

	for sTime := range cronJobs {
		s.Add(sTime, func() {
			js.ExecuteJobs(cronJobs[sTime])
		})
	}

	if config.GetDeleteRunsAfterDays() > 0 {
		s.Add("0 0 * * *", func() {
			queries.DeleteRuns(context.Background(), time.Now().AddDate(0, 0, -int(config.GetDeleteRunsAfterDays())).UnixMilli())
		})
	}

	js.Events = events.New(func(streamID string, sub *sse.Subscriber) {
		all := js.ListJobs()
		js.Events.SendEvent(js.IsIdle(), nil, &all)
	})

	return js, nil
}

type JobService struct {
	Queries   *jobs.Queries
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

func (js *JobService) ExecuteJobs(jobs map[string]config.Job) {
	if len(jobs) == 0 {
		jobs = config.GetJobs()
	}
	for name, job := range jobs {
		js.ExecuteJob(name, &job)
	}
}

func (js *JobService) ExecuteJob(name string, job *config.Job) {
	ctx := context.Background()

	runView := js.startRun(ctx, name)

	envs := config.GetEnvsByJobName(name)
	for key, value := range envs {
		os.Setenv(key, os.ExpandEnv(value))
	}

	// js.writeLog(ctx, name, runView.ID, Debug, fmt.Sprintf("Setting environment variables:\n\t%s", strings.Join(keys, "\n\t")))

	for _, command := range config.GetCommandsByJobName(name) {
		severity := Debug
		// js.writeLog(ctx, dbJob, runView.ID, Debug, fmt.Sprintf("Executing command: %s", command.Command))
		out, err := commands.ExecuteCommand(command)
		severity = Info
		if err != nil {
			severity = Error
		}
		fmt.Println(severity, out)
		// js.writeLog(ctx, dbJob, runView.ID, severity, out)
		if err != nil {
			runView.StatusID = Stopped.Int64()
			break
		}
	}

	for key := range envs {
		os.Unsetenv(key)
	}

	js.endRun(ctx, name, runView)
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

// func (js *JobService) refreshLogs(jobView *jobs.JobsView, run *jobs.Run, newLog *jobs.Log) {
// 	amount := len(jobView.Runs)
// 	// most likely the last run
// 	for i := amount - 1; i >= 0; i-- {
// 		if jobView.Runs[i].ID == run.ID {
// 			createdAtSeconds := newLog.CreatedAt / 1000
// 			t := time.Unix(createdAtSeconds, 0).Local()
// 			formattedTime := t.Format(DATE_FORMAT)

// 			jobView.Runs[i].Logs = append(jobView.Runs[i].Logs, jobs.ListLogsByRunIDRow{
// 				CreatedAt:     newLog.CreatedAt,
// 				RunID:         newLog.RunID,
// 				SeverityID:    newLog.SeverityID,
// 				Message:       newLog.Message,
// 				CreatedAtTime: formattedTime,
// 			})

// 			break
// 		}
// 	}
// }

func (js *JobService) startRun(ctx context.Context, jobName string) *jobs.RunsView {
	run, _ := js.Queries.CreateRun(ctx, jobs.CreateRunParams{
		JobID:     jobName,
		StatusID:  Running.Int64(),
		StartTime: time.Now().UnixMilli(),
	})

	runView := &jobs.RunsView{
		ID:           run.ID,
		JobID:        jobName,
		StatusID:     run.StatusID,
		StartTime:    run.StartTime,
		EndTime:      run.EndTime,
		FmtStartTime: formatTime(run.StartTime),
		Logs:         nil,
	}
	// dbJob.Runs = append(dbJob.Runs, *runView)

	// js.Events.SendEvent(true, dbJob, nil)
	// prepare run to be finished if no error is set
	runView.StatusID = Finished.Int64()
	return runView
}

func (js *JobService) endRun(ctx context.Context, name string, runView *jobs.RunsView) {
	// run, _ := js.Queries.UpdateRun(ctx, jobs.UpdateRunParams{
	// 	StatusID: runView.StatusID,
	// 	EndTime:  sql.NullInt64{Int64: time.Now().UnixMilli(), Valid: true},
	// 	ID:       runView.ID,
	// })

	// amount := len(dbJob.Runs)
	// // most likely the last run
	// for i := amount - 1; i >= 0; i-- {
	// 	if dbJob.Runs[i].ID == run.ID {
	// 		dbJob.Runs[i].FmtEndTime.String = formatTime(run.EndTime.Int64)
	// 		dbJob.Runs[i].FmtEndTime.Valid = true
	// 		dbJob.Runs[i].Duration.Int64 = run.EndTime.Int64 - run.StartTime
	// 		dbJob.Runs[i].Duration.Valid = true
	// 		dbJob.Runs[i].StatusID = run.StatusID
	// 		dbJob.Runs[i].EndTime = run.EndTime
	// 		break
	// 	}
	// }

	// js.Events.SendEvent(true, dbJob, nil)
}

// func (js *JobService) writeLog(ctx context.Context, dbJob *jobs.JobsView, runId int64, severity Severity, message string) {
// 	newLog, _ := js.Queries.CreateLog(ctx, jobs.CreateLogParams{
// 		CreatedAt:  time.Now().UnixMilli(),
// 		RunID:      runId,
// 		SeverityID: int64(severity),
// 		Message:    message,
// 	})

// 	js.refreshLogs(dbJob, &jobs.Run{ID: runId}, &newLog)
// 	js.Events.SendEvent(false, dbJob, nil)
// }
