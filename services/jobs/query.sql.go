// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package jobs

import (
	"context"
)

const createJob = `-- name: CreateJob :one
INSERT INTO
    jobs (Name, backup_schedule_id)
VALUES
    (?, ?) RETURNING id, name, backup_schedule_id
`

type CreateJobParams struct {
	Name             string `json:"name"`
	BackupScheduleID int64  `json:"backup_schedule_id"`
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, createJob, arg.Name, arg.BackupScheduleID)
	var i Job
	err := row.Scan(&i.ID, &i.Name, &i.BackupScheduleID)
	return i, err
}

const getJob = `-- name: GetJob :one
SELECT
    id, name, backup_schedule_id
FROM
    jobs
WHERE
    id = ?
`

func (q *Queries) GetJob(ctx context.Context, id int64) (Job, error) {
	row := q.db.QueryRowContext(ctx, getJob, id)
	var i Job
	err := row.Scan(&i.ID, &i.Name, &i.BackupScheduleID)
	return i, err
}

const listJobs = `-- name: ListJobs :many
SELECT
    id, name, backup_schedule_id
FROM
    jobs
ORDER BY
    name
`

func (q *Queries) ListJobs(ctx context.Context) ([]Job, error) {
	rows, err := q.db.QueryContext(ctx, listJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(&i.ID, &i.Name, &i.BackupScheduleID); err != nil {
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

const updateJob = `-- name: UpdateJob :exec
UPDATE jobs
set
    name = ?,
    backup_schedule_id = ?
WHERE
    id = ?
`

type UpdateJobParams struct {
	Name             string `json:"name"`
	BackupScheduleID int64  `json:"backup_schedule_id"`
	ID               int64  `json:"id"`
}

func (q *Queries) UpdateJob(ctx context.Context, arg UpdateJobParams) error {
	_, err := q.db.ExecContext(ctx, updateJob, arg.Name, arg.BackupScheduleID, arg.ID)
	return err
}
