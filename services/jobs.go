package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"

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

var (
	lastTimestamp int64
	mu            sync.Mutex
)

func generateUniqueTimestamp() int64 {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UnixMilli()
	if now <= lastTimestamp {
		lastTimestamp++
	} else {
		lastTimestamp = now
	}
	return lastTimestamp
}

type JobView struct {
	Name string    `json:"name"`
	Cron string    `json:"cron"`
	Runs []RunView `json:"runs"`
}

type RunView struct {
	ID        int64                      `json:"id"`
	JobName   string                     `json:"job_name"`
	StatusID  int64                      `json:"status_id"`
	StartTime string                     `json:"start_time"`
	EndTime   string                     `json:"end_time"`
	Duration  string                     `json:"duration"`
	Logs      []jobs.ListLogsByRunIDsRow `json:"logs"`
}

func NewJobService() (*JobService, error) {
	s := scheduler.New()
	queries, err := setupSQLite()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var cronJobs = make(map[string][]config.Job)
	js := &JobService{Queries: queries, Scheduler: s}

	jobs := config.GetJobs()
	for _, job := range jobs {
		cron := config.GetJobsCron(&job)
		cronJobs[cron] = append(cronJobs[cron], job)
	}

	for sTime := range cronJobs {
		s.Add(sTime, func() {
			js.ExecuteJobs(cronJobs[sTime])
		})
	}

	if config.GetDeleteRunsAfterDays() > 0 {
		s.Add("0 0 * * *", func() {
			queries.DeleteOldRuns(context.Background(), time.Now().AddDate(0, 0, -int(config.GetDeleteRunsAfterDays())).UnixMilli())
		})
	}

	deleteOrphanedRuns(queries)
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Config file changed", "path", e.Name)
		deleteOrphanedRuns(queries)
		js.Events.SendJobEvent(js.IsIdle(), nil, js.ListJobs())
	})
	viper.WatchConfig()

	return js, nil
}

type JobService struct {
	Queries   *jobs.Queries
	Scheduler *scheduler.Scheduler
	Events    *events.Event
}

func (js *JobService) SetEvents(e *events.Event) {
	js.Events = e
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

func deleteOrphanedRuns(queries *jobs.Queries) {
	names := []sql.NullString{}
	j := config.GetJobs()
	for _, job := range j {
		names = append(names, sql.NullString{String: strings.ToLower(job.Name), Valid: true})
	}
	queries.DeleteObsoleteRuns(context.Background(), names)
}

func (js *JobService) ExecuteJobs(jobs []config.Job) {
	if len(jobs) == 0 {
		jobs = config.GetJobs()
	}
	for _, job := range jobs {
		js.ExecuteJob(&job)
	}
}

func (js *JobService) ExecuteJob(job *config.Job) {
	ctx := context.Background()

	run, err := js.startRun(ctx, job.Name)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Key storage for log
	keys := []string{}
	envs := config.GetEnvsByJobName(job.Name)
	for key, value := range envs {
		os.Setenv(key, os.ExpandEnv(value))
		keys = append(keys, key)
	}

	js.writeLog(ctx, run, Debug, fmt.Sprintf("Setting environment variables: %s", strings.Join(keys, ", ")))

	for _, command := range config.GetCommandsByJobName(job.Name) {
		severity := Debug
		js.writeLog(ctx, run, Debug, fmt.Sprintf("Executing command: %s", command))
		out, err := commands.ExecuteCommand(command)
		severity = Info
		if err != nil {
			severity = Error
		}
		js.writeLog(ctx, run, severity, out)
		if err != nil {
			run.StatusID = Stopped.Int64()
		} else {
			run.StatusID = Finished.Int64()
		}
	}

	for key := range envs {
		os.Unsetenv(key)
	}

	js.endRun(ctx, run)
}

func (js *JobService) ListJobs() []JobView {
	jobs := config.GetJobs()

	runs, err := js.Queries.GetThreeRunsPerJobName(context.Background())
	if err != nil {
		slog.Error(err.Error())
		return []JobView{}
	}

	runsByJob := make(map[string][]RunView)
	for _, run := range runs {
		endTime := ""
		var duration time.Duration
		if run.EndTime.Valid {
			endTime = formatTime(run.EndTime.Int64)
			duration = time.Duration(run.EndTime.Int64-run.StartTime) * time.Millisecond
		}
		runsByJob[run.JobName] = append(runsByJob[run.JobName],
			RunView{
				ID:        run.ID,
				StatusID:  run.StatusID,
				StartTime: formatTime(run.StartTime),
				EndTime:   endTime,
				Duration:  duration.Truncate(time.Second).String(),
			})
	}

	result := make([]JobView, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, JobView{
			Name: job.Name,
			Cron: config.GetJobsCron(&job),
			Runs: runsByJob[job.Name],
		})
	}

	return result
}

func (js *JobService) ListRuns(name string, limit int64) ([]RunView, error) {
	normalized := sql.NullString{String: strings.ToLower(name), Valid: true}

	runs, err := js.Queries.GetRuns(context.Background(), jobs.GetRunsParams{JobNameNormalized: normalized, Limit: limit})
	if err != nil {
		return nil, fmt.Errorf("failed to get runs for job %s: %w", name, err)
	}

	if len(runs) == 0 {
		return []RunView{}, nil
	}

	runIDs := make([]int64, 0, len(runs))
	for _, run := range runs {
		runIDs = append(runIDs, run.ID)
	}

	allLogs, err := js.Queries.ListLogsByRunIDs(context.Background(), runIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for runs: %w", err)
	}

	logsByRun := make(map[int64][]jobs.ListLogsByRunIDsRow)
	for _, log := range allLogs {
		logsByRun[log.RunID] = append(logsByRun[log.RunID], log)
	}

	result := make([]RunView, 0, len(runs))
	for _, run := range runs {
		endTime := ""
		var duration time.Duration
		if run.EndTime.Valid {
			endTime = formatTime(run.EndTime.Int64)
			duration = time.Duration(run.EndTime.Int64-run.StartTime) * time.Millisecond
		}

		result = append(result, RunView{
			ID:        run.ID,
			JobName:   run.JobName,
			StatusID:  run.StatusID,
			StartTime: formatTime(run.StartTime),
			EndTime:   endTime,
			Duration:  duration.Truncate(time.Second).String(),
			Logs:      logsByRun[run.ID],
		})
	}

	return result, nil
}

func (js *JobService) startRun(ctx context.Context, jobName string) (*jobs.Run, error) {
	run, err := js.Queries.CreateRun(ctx, jobs.CreateRunParams{
		JobName:   jobName,
		StatusID:  Running.Int64(),
		StartTime: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, err
	}
	js.Events.SendJobEvent(true, js.getLatestRun(ctx, &run), nil)
	return &run, nil
}

func (js *JobService) endRun(ctx context.Context, run *jobs.Run) {
	_, err := js.Queries.UpdateRun(ctx, jobs.UpdateRunParams{
		StatusID: run.StatusID,
		EndTime:  sql.NullInt64{Int64: time.Now().UnixMilli(), Valid: true},
		ID:       run.ID,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	js.Events.SendJobEvent(true, js.getLatestRun(ctx, run), nil)
}

func (js *JobService) writeLog(ctx context.Context, run *jobs.Run, severity Severity, message string) {
	_, err := js.Queries.CreateLog(ctx, jobs.CreateLogParams{
		CreatedAt:  generateUniqueTimestamp(),
		RunID:      run.ID,
		SeverityID: int64(severity),
		Message:    message,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	js.Events.SendJobEvent(false, js.getLatestRun(ctx, run), nil)
}

func (js *JobService) getLatestRun(ctx context.Context, run *jobs.Run) *RunView {
	runs, err := js.Queries.GetRuns(ctx, jobs.GetRunsParams{JobNameNormalized: run.JobNameNormalized, Limit: 1})
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	if len(runs) == 0 {
		slog.Warn("No new run found")
		return nil
	}
	r := runs[0]
	logs, err := js.Queries.ListLogsByRunIDs(context.Background(), []int64{r.ID})
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	endTime := ""
	var duration time.Duration
	if r.EndTime.Valid {
		endTime = formatTime(r.EndTime.Int64)
		duration = time.Duration(r.EndTime.Int64-r.StartTime) * time.Millisecond
	}
	return &RunView{
		ID:        r.ID,
		JobName:   r.JobName,
		StatusID:  r.StatusID,
		StartTime: formatTime(r.StartTime),
		EndTime:   endTime,
		Duration:  duration.Truncate(time.Second).String(),
		Logs:      logs,
	}
}
