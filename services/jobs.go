package services

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"

	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/internal/commands"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

//go:embed jobs.sql
var ddl string

func NewJobService(dbName string, config *config.Config) (*JobService, error) {
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

	return &JobService{Queries: queries, Config: config}, nil
}

type JobService struct {
	Queries *jobs.Queries
	Config  *config.Config
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
	dbJobs, err := queries.ListJobs(ctx)
	if err != nil {
		return err
	}

	existingJobs := make(map[string]bool)
	for _, job := range dbJobs {
		existingJobs[job] = true
	}

	for _, job := range config.Jobs {
		queries.CreateJob(ctx, job.Name)
		delete(existingJobs, job.Name)
	}

	for name := range existingJobs {
		err := queries.DeleteJob(ctx, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (js *JobService) GetQueries() *jobs.Queries {
	return js.Queries
}

func (js *JobService) ExecuteJobs() {
	for i := 0; i < len(js.Config.Jobs); i++ {
		js.ExecuteJob(i)
	}
}

func (js *JobService) ExecuteJob(id int) {
	job := js.Config.Jobs[id]
	ctx := context.Background()

	run, _ := js.Queries.CreateRun(ctx, jobs.CreateRunParams{
		Job:      job.Name,
		StatusID: int64(Running),
	})
	status := Finished

	for _, command := range job.Envs {
		js.Queries.CreateLog(ctx, jobs.CreateLogParams{
			RunID:      run.ID,
			SeverityID: int64(Debug),
			Message:    fmt.Sprintf("Setting environment variable: \"%s\"", command.Key),
		})
		os.Setenv(command.Key, command.Value)
	}

	for _, command := range job.Commands {
		program, args, err := commands.PrepareCommand(command.Command)
		if err != nil {
		}
		js.Queries.CreateLog(ctx, jobs.CreateLogParams{
			RunID:      run.ID,
			SeverityID: int64(Debug),
			Message:    fmt.Sprintf("Executing command: \"%s\" - \"%s\"", program, args),
		})
		out, err := commands.ExecuteCommand(program, args)
		severity := Info
		if err != nil {
			severity = Error
		}
		js.Queries.CreateLog(ctx, jobs.CreateLogParams{
			RunID:      run.ID,
			SeverityID: int64(severity),
			Message:    out,
		})
		if err != nil {
			status = Stopped
			break
		}
		// Wait for 2 seconds to make sure the command is finished
		time.Sleep(2 * time.Second)
	}

	js.Queries.UpdateRun(ctx, jobs.UpdateRunParams{
		StatusID: int64(status),
		EndTime:  sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ID:       run.ID,
	})
}
