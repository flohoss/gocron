package services

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"

	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

//go:embed jobs.sql
var ddl string

func NewJobService(dbName string) (*jobs.Queries, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", dbName+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}
	log.Println("🚀 Connected Successfully to the Database")

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return jobs.New(db), nil
}
