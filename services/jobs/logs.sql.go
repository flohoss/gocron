// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: logs.sql

package jobs

import (
	"context"
)

const createLog = `-- name: CreateLog :one
INSERT INTO
    logs (created_at, run_id, severity_id, message)
VALUES
    (?, ?, ?, ?) RETURNING created_at, run_id, severity_id, message
`

type CreateLogParams struct {
	CreatedAt  int64  `json:"created_at"`
	RunID      int64  `json:"run_id"`
	SeverityID int64  `json:"severity_id"`
	Message    string `json:"message"`
}

func (q *Queries) CreateLog(ctx context.Context, arg CreateLogParams) (Log, error) {
	row := q.db.QueryRowContext(ctx, createLog,
		arg.CreatedAt,
		arg.RunID,
		arg.SeverityID,
		arg.Message,
	)
	var i Log
	err := row.Scan(
		&i.CreatedAt,
		&i.RunID,
		&i.SeverityID,
		&i.Message,
	)
	return i, err
}

const listLogsByRunID = `-- name: ListLogsByRunID :many
SELECT
    created_at, run_id, severity_id, message,
    STRFTIME(
        '%H:%M:%S %Y-%m-%d',
        created_at / 1000,
        'unixepoch',
        'localtime'
    ) AS created_at_time
FROM
    logs
WHERE
    run_id = ?
ORDER BY
    created_at
`

type ListLogsByRunIDRow struct {
	CreatedAt     int64       `json:"created_at"`
	RunID         int64       `json:"run_id"`
	SeverityID    int64       `json:"severity_id"`
	Message       string      `json:"message"`
	CreatedAtTime interface{} `json:"created_at_time"`
}

func (q *Queries) ListLogsByRunID(ctx context.Context, runID int64) ([]ListLogsByRunIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listLogsByRunID, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLogsByRunIDRow
	for rows.Next() {
		var i ListLogsByRunIDRow
		if err := rows.Scan(
			&i.CreatedAt,
			&i.RunID,
			&i.SeverityID,
			&i.Message,
			&i.CreatedAtTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
