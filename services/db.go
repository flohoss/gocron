package services

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
)

//go:embed jobs.sql
var ddl string

func setupSQLite() (*jobs.Queries, error) {
	ctx := context.Background()

	dbLocation := config.GetDBLocation()
	if err := os.MkdirAll(dbLocation, os.ModePerm); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbLocation, config.GetDBName())
	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	queries := jobs.New(db)

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	// TODO: Replace this ad-hoc migration with a versioned migration system.
	// Create services/migrations/ with numbered files (e.g. 001_initial.sql,
	// 002_add_job_slug.sql) and a schema_version table to track which migrations
	// have been applied. Each migration must be reviewed and updated on every
	// release that touches the DB schema.
	if err := migrateSchema(ctx, db); err != nil {
		return nil, err
	}

	return queries, nil
}

func migrateSchema(ctx context.Context, db *sql.DB) error {
	var colName string
	err := db.QueryRowContext(ctx, "SELECT name FROM pragma_table_info('runs') WHERE name = 'job_name_normalized'").Scan(&colName)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	slog.Info("Migrating runs table from job_name_normalized to job_slug")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var hasSlug bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pragma_table_info('runs') WHERE name = 'job_slug')").Scan(&hasSlug)
	if err != nil {
		return err
	}

	if !hasSlug {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE runs ADD COLUMN job_slug TEXT NOT NULL DEFAULT ''`); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx, `UPDATE runs SET job_slug = LOWER(TRIM(job_name)) WHERE job_slug = ''`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DROP INDEX IF EXISTS idx_runs_job_name_normalized_start_time`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_runs_job_slug_start_time ON runs (job_slug, start_time DESC)`); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	slog.Info("Migration complete: job_name_normalized replaced with job_slug")
	return nil
}
