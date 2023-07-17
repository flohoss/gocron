package database

import (
	"context"
	"database/sql"
	_ "embed"
	"os"

	_ "github.com/glebarez/sqlite"
)

const Storage = "storage/"

func init() {
	os.Mkdir(Storage, os.ModePerm)
}

//go:embed sql/schema.sql
var ddl string

func MigrateDatabase() (*Queries, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", Storage+"db.sqlite3?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return New(db), nil
}
