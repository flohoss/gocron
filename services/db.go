package services

import (
	"context"
	"database/sql"
	_ "embed"
	"path/filepath"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/services/jobs"
)

//go:embed jobs.sql
var ddl string

func setupSQLite() (*jobs.Queries, error) {
	ctx := context.Background()

	dbPath := filepath.Join(config.GetConfigFolderPath(), "db.sqlite")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	queries := jobs.New(db)

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return queries, nil
}
