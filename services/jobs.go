package services

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/glebarez/go-sqlite"

	"gitlab.unjx.de/flohoss/gobackup/config"
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

	err = createUpdateOrDeleteJob(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	return &JobService{Queries: queries}, nil
}

type JobService struct {
	Queries *jobs.Queries
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
