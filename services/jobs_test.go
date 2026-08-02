package services

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/internal/events"
	"github.com/flohoss/gocron/services/jobs"
	"github.com/r3labs/sse/v2"
	"github.com/spf13/viper"
)

func TestFormatTime_FormatsUnixMillis(t *testing.T) {
	ms := int64(1753900800000)
	got := formatTime(ms)
	want := time.Unix(ms/1000, 0).Local().Format(DATE_FORMAT)
	if got != want {
		t.Fatalf("formatTime(%d) = %q, want %q", ms, got, want)
	}
}

func TestFormatTime_HandlesZero(t *testing.T) {
	got := formatTime(0)
	if got != "1970-01-01 01:00:00" && got != "1970-01-01 00:00:00" {
		t.Fatalf("unexpected format for zero time: %q", got)
	}
}

func TestGenerateUniqueTimestamp_AlwaysIncrements(t *testing.T) {
	lastTimestamp = 0

	first := generateUniqueTimestamp()
	second := generateUniqueTimestamp()

	if second <= first {
		t.Fatalf("expected second > first, got first=%d second=%d", first, second)
	}
}

func TestGenerateUniqueTimestamp_IncrementsWhenSameMillisecond(t *testing.T) {
	lastTimestamp = time.Now().UnixMilli()

	a := generateUniqueTimestamp()
	b := generateUniqueTimestamp()

	if b != a+1 {
		t.Fatalf("expected b=a+1 when same ms, got a=%d b=%d", a, b)
	}
}

func TestGenerateUniqueTimestamp_ConcurrentCallsAreUnique(t *testing.T) {
	lastTimestamp = 0

	var wg sync.WaitGroup
	results := make(chan int64, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- generateUniqueTimestamp()
		}()
	}

	wg.Wait()
	close(results)

	seen := make(map[int64]bool, 100)
	for ts := range results {
		if seen[ts] {
			t.Fatalf("duplicate timestamp generated: %d", ts)
		}
		seen[ts] = true
	}
}

func TestExecuteJobs_SkipsWhenNotIdle(t *testing.T) {
	js := &JobService{
		jobCtx:    context.Background(),
		jobCancel: func() {},
	}

	js.jobMu.Lock()
	js.jobRunning = true
	js.jobMu.Unlock()

	called := false
	js.ExecuteJobs(nil)
	if called {
		t.Fatal("expected ExecuteJobs to skip when job is already running")
	}
}

func setupTestDB(t *testing.T) *jobs.Queries {
	t.Helper()

	tmpDir := t.TempDir()

	v := viper.New()
	v.Set("log_level", "info")
	v.Set("time_zone", "UTC")
	v.Set("delete_runs_after_days", 0)
	v.Set("db.location", tmpDir)
	v.Set("db.name", "test.sqlite")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("terminal.allow_all_commands", true)

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	queries, err := setupSQLite()
	if err != nil {
		t.Fatalf("failed to setup SQLite: %v", err)
	}

	t.Cleanup(func() {
		queries.StopRunning(context.Background())
	})

	return queries
}

func TestStopRunning_MarksRunningAsCanceled(t *testing.T) {
	queries := setupTestDB(t)

	run, err := queries.CreateRun(context.Background(), jobs.CreateRunParams{
		JobName:   "test job",
		JobSlug:   "test-job",
		StatusID:  Running.Int64(),
		StartTime: time.Now().UnixMilli(),
	})
	if err != nil {
		t.Fatalf("failed to create run: %v", err)
	}

	if run.StatusID != Running.Int64() {
		t.Fatalf("expected running status, got %d", run.StatusID)
	}

	if err := queries.StopRunning(context.Background()); err != nil {
		t.Fatalf("failed to stop running: %v", err)
	}

	runs, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: run.JobSlug,
		Limit:   1,
	})
	if err != nil {
		t.Fatalf("failed to get runs: %v", err)
	}
	if len(runs) == 0 {
		t.Fatal("expected at least one run")
	}
	if runs[0].StatusID != Canceled.Int64() {
		t.Fatalf("expected canceled status (%d), got %d", Canceled.Int64(), runs[0].StatusID)
	}
	if !runs[0].EndTime.Valid {
		t.Fatal("expected end_time to be set")
	}
}

func TestIsIdle_TrueWhenNoRunningJobs(t *testing.T) {
	queries := setupTestDB(t)

	js := &JobService{Queries: queries}
	if !js.IsIdle() {
		t.Fatal("expected idle when no running jobs")
	}
}

func TestIsIdle_FalseWhenJobRunning(t *testing.T) {
	queries := setupTestDB(t)

	_, err := queries.CreateRun(context.Background(), jobs.CreateRunParams{
		JobName:   "running job",
		JobSlug:   "running-job",
		StatusID:  Running.Int64(),
		StartTime: time.Now().UnixMilli(),
	})
	if err != nil {
		t.Fatalf("failed to create run: %v", err)
	}

	js := &JobService{Queries: queries}
	if js.IsIdle() {
		t.Fatal("expected not idle when a job is running")
	}
}

func TestExecuteJob_Success(t *testing.T) {
	queries := setupTestDB(t)
	ctx, cancel := context.WithCancel(context.Background())

	js := &JobService{
		Queries:   queries,
		Events:    events.New(func(string, *sse.Subscriber) {}),
		jobCtx:    ctx,
		jobCancel: cancel,
	}

	job := &config.Job{
		Name:     "test success",
		Slug:     "test-success",
		Commands: []string{"echo hello"},
	}

	js.ExecuteJob(job)

	runs, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: "test-success",
		Limit:   1,
	})
	if err != nil {
		t.Fatalf("failed to get runs: %v", err)
	}
	if len(runs) == 0 {
		t.Fatal("expected at least one run")
	}
	if runs[0].StatusID != Finished.Int64() {
		t.Fatalf("expected finished status, got %d", runs[0].StatusID)
	}
}

func TestExecuteJob_Failure(t *testing.T) {
	queries := setupTestDB(t)
	ctx, cancel := context.WithCancel(context.Background())

	js := &JobService{
		Queries:   queries,
		Events:    events.New(func(string, *sse.Subscriber) {}),
		jobCtx:    ctx,
		jobCancel: cancel,
	}

	job := &config.Job{
		Name:     "test failure",
		Slug:     "test-failure",
		Commands: []string{"exit 1"},
	}

	js.ExecuteJob(job)

	runs, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: "test-failure",
		Limit:   1,
	})
	if err != nil {
		t.Fatalf("failed to get runs: %v", err)
	}
	if len(runs) == 0 {
		t.Fatal("expected at least one run")
	}
	if runs[0].StatusID != Stopped.Int64() {
		t.Fatalf("expected stopped status, got %d", runs[0].StatusID)
	}
}

func TestExecuteJob_Canceled(t *testing.T) {
	queries := setupTestDB(t)
	ctx, cancel := context.WithCancel(context.Background())

	js := &JobService{
		Queries:   queries,
		Events:    events.New(func(string, *sse.Subscriber) {}),
		jobCtx:    ctx,
		jobCancel: cancel,
	}

	job := &config.Job{
		Name:     "test canceled",
		Slug:     "test-canceled",
		Commands: []string{"sleep 30"},
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	js.ExecuteJob(job)

	runs, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: "test-canceled",
		Limit:   1,
	})
	if err != nil {
		t.Fatalf("failed to get runs: %v", err)
	}
	if len(runs) == 0 {
		t.Fatal("expected at least one run")
	}
	if runs[0].StatusID != Canceled.Int64() {
		t.Fatalf("expected canceled status, got %d", runs[0].StatusID)
	}
}

func TestExecuteJob_Timeout(t *testing.T) {
	queries := setupTestDB(t)
	ctx, cancel := context.WithCancel(context.Background())

	js := &JobService{
		Queries:   queries,
		Events:    events.New(func(string, *sse.Subscriber) {}),
		jobCtx:    ctx,
		jobCancel: cancel,
	}

	job := &config.Job{
		Name:     "test timeout",
		Slug:     "test-timeout",
		Timeout:  1 * time.Second,
		Commands: []string{"sleep 30"},
	}

	js.ExecuteJob(job)

	runs, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: "test-timeout",
		Limit:   1,
	})
	if err != nil {
		t.Fatalf("failed to get runs: %v", err)
	}
	if len(runs) == 0 {
		t.Fatal("expected at least one run")
	}
	if runs[0].StatusID != Stopped.Int64() {
		t.Fatalf("expected stopped status (timeout), got %d", runs[0].StatusID)
	}
}

func TestDeleteOldRuns(t *testing.T) {
	queries := setupTestDB(t)

	oldRun, err := queries.CreateRun(context.Background(), jobs.CreateRunParams{
		JobName:   "old job",
		JobSlug:   "old-job",
		StatusID:  Finished.Int64(),
		StartTime: time.Now().AddDate(0, 0, -30).UnixMilli(),
	})
	if err != nil {
		t.Fatalf("failed to create old run: %v", err)
	}

	_, err = queries.CreateRun(context.Background(), jobs.CreateRunParams{
		JobName:   "recent job",
		JobSlug:   "recent-job",
		StatusID:  Finished.Int64(),
		StartTime: time.Now().UnixMilli(),
	})
	if err != nil {
		t.Fatalf("failed to create recent run: %v", err)
	}

	err = queries.DeleteOldRuns(context.Background(), time.Now().AddDate(0, 0, -7).UnixMilli())
	if err != nil {
		t.Fatalf("failed to delete old runs: %v", err)
	}

	oldRuns, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: oldRun.JobSlug,
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("failed to get old runs: %v", err)
	}
	if len(oldRuns) != 0 {
		t.Fatalf("expected old run to be deleted, got %d runs", len(oldRuns))
	}
}

func TestDeleteObsoleteRuns(t *testing.T) {
	queries := setupTestDB(t)

	obsoleteRun, err := queries.CreateRun(context.Background(), jobs.CreateRunParams{
		JobName:   "obsolete job",
		JobSlug:   "obsolete-job",
		StatusID:  Finished.Int64(),
		StartTime: time.Now().UnixMilli(),
	})
	if err != nil {
		t.Fatalf("failed to create obsolete run: %v", err)
	}

	err = queries.DeleteObsoleteRuns(context.Background(), []string{
		"active-job",
	})
	if err != nil {
		t.Fatalf("failed to delete obsolete runs: %v", err)
	}

	obsoleteRuns, err := queries.GetRuns(context.Background(), jobs.GetRunsParams{
		JobSlug: obsoleteRun.JobSlug,
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("failed to get obsolete runs: %v", err)
	}
	if len(obsoleteRuns) != 0 {
		t.Fatalf("expected obsolete run to be deleted, got %d runs", len(obsoleteRuns))
	}
}

func TestSetupSQLite_CreatesDatabase(t *testing.T) {
	tmpDir := t.TempDir()

	v := viper.New()
	v.Set("log_level", "info")
	v.Set("time_zone", "UTC")
	v.Set("delete_runs_after_days", 0)
	v.Set("db.location", tmpDir)
	v.Set("db.name", "created.sqlite")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("terminal.allow_all_commands", true)

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	queries, err := setupSQLite()
	if err != nil {
		t.Fatalf("failed to setup SQLite: %v", err)
	}
	if queries == nil {
		t.Fatal("expected non-nil queries")
	}

	dbPath := filepath.Join(tmpDir, "created.sqlite")
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected database file at %s: %v", dbPath, err)
	}
}

func TestGetDBLocation_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()

	config.SetConfigFilePath(filepath.Join(tmpDir, "config.yaml"))

	v := viper.New()
	v.Set("db.location", ".")
	v.Set("db.name", "test.db")
	v.Set("server.address", "127.0.0.1")
	v.Set("server.port", 8156)
	v.Set("log_level", "info")
	v.Set("time_zone", "UTC")
	v.Set("delete_runs_after_days", 0)
	v.Set("terminal.allow_all_commands", true)

	if err := config.ValidateAndLoadConfig(v); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	loc := config.GetDBLocation()
	if loc != tmpDir {
		t.Fatalf("expected db location %s, got %s", tmpDir, loc)
	}

	name := config.GetDBName()
	if name != "test.db" {
		t.Fatalf("expected db name test.db, got %s", name)
	}
}
