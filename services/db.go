package services

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/flohoss/gocron/config"
	"github.com/flohoss/gocron/internal/buildinfo"
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

	if err := purgeIfOutdated(ctx, dbPath); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		slog.Warn("Failed to apply DDL, purging database and retrying", "path", dbPath, "error", err)
		if err := db.Close(); err != nil {
			return nil, err
		}
		if err := os.Remove(dbPath); err != nil {
			return nil, err
		}
		db, err = sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
		if err != nil {
			return nil, err
		}
		if _, err := db.ExecContext(ctx, ddl); err != nil {
			return nil, err
		}
	}

	queries := jobs.New(db)
	if err := queries.SetSchemaVersion(ctx, buildinfo.Version); err != nil {
		return nil, err
	}

	return queries, nil
}

func purgeIfOutdated(ctx context.Context, dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return err
	}
	defer db.Close()

	version, err := jobs.New(db).GetSchemaVersion(ctx)
	if err == nil && version == buildinfo.Version {
		return nil
	}

	if err != nil {
		slog.Warn("Database schema is unversioned, purging existing data", "path", dbPath, "error", err)
	} else {
		slog.Warn("Database schema is outdated, purging existing data", "path", dbPath, "version", version)
	}
	return os.Remove(dbPath)
}
