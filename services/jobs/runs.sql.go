// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: runs.sql

package jobs

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const createRun = `-- name: CreateRun :one
INSERT INTO
    runs (job_id, status_id, start_time)
VALUES
    (?, ?, ?) RETURNING id, job_id, status_id, start_time, end_time
`

type CreateRunParams struct {
	JobID     string `json:"job_id"`
	StatusID  int64  `json:"status_id"`
	StartTime int64  `json:"start_time"`
}

func (q *Queries) CreateRun(ctx context.Context, arg CreateRunParams) (Run, error) {
	row := q.db.QueryRowContext(ctx, createRun, arg.JobID, arg.StatusID, arg.StartTime)
	var i Run
	err := row.Scan(
		&i.ID,
		&i.JobID,
		&i.StatusID,
		&i.StartTime,
		&i.EndTime,
	)
	return i, err
}

const getRunsView = `-- name: GetRunsView :many
SELECT
    id, job_id, status_id, start_time, end_time, fmt_start_time, fmt_end_time, duration, logs
FROM
    (
        SELECT
            id, job_id, status_id, start_time, end_time, fmt_start_time, fmt_end_time, duration, logs
        FROM
            runs_view
        WHERE
            job_id = ?
        ORDER BY
            start_time DESC
        LIMIT
            5
    ) subquery
ORDER BY
    start_time ASC
`

func (q *Queries) GetRunsView(ctx context.Context, jobID string) ([]RunsView, error) {
	rows, err := q.db.QueryContext(ctx, getRunsView, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RunsView
	for rows.Next() {
		var i RunsView
		if err := rows.Scan(
			&i.ID,
			&i.JobID,
			&i.StatusID,
			&i.StartTime,
			&i.EndTime,
			&i.FmtStartTime,
			&i.FmtEndTime,
			&i.Duration,
			pq.Array(&i.Logs),
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

const isIdle = `-- name: IsIdle :one
SELECT
    CASE
        WHEN status_id = 1 THEN FALSE
        ELSE TRUE
    END AS is_idle
FROM
    runs
UNION ALL
SELECT
    TRUE
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            runs
    )
LIMIT
    1
`

func (q *Queries) IsIdle(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, isIdle)
	var is_idle int64
	err := row.Scan(&is_idle)
	return is_idle, err
}

const updateRun = `-- name: UpdateRun :exec
UPDATE runs
SET
    status_id = ?,
    end_time = ?
WHERE
    id = ?
`

type UpdateRunParams struct {
	StatusID int64         `json:"status_id"`
	EndTime  sql.NullInt64 `json:"end_time"`
	ID       int64         `json:"id"`
}

func (q *Queries) UpdateRun(ctx context.Context, arg UpdateRunParams) error {
	_, err := q.db.ExecContext(ctx, updateRun, arg.StatusID, arg.EndTime, arg.ID)
	return err
}
