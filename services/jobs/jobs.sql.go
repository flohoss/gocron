// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: jobs.sql

package jobs

import (
	"context"
)

const createJob = `-- name: CreateJob :one
INSERT INTO
    jobs (id, name, cron)
VALUES
    (?, ?, ?) RETURNING id, name, cron
`

type CreateJobParams struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Cron string `json:"cron"`
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, createJob, arg.ID, arg.Name, arg.Cron)
	var i Job
	err := row.Scan(&i.ID, &i.Name, &i.Cron)
	return i, err
}

const deleteJob = `-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    id = ?
`

func (q *Queries) DeleteJob(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteJob, id)
	return err
}

const getJob = `-- name: GetJob :one
SELECT
    id, name, cron
FROM
    jobs
WHERE
    id = ?
`

func (q *Queries) GetJob(ctx context.Context, id string) (Job, error) {
	row := q.db.QueryRowContext(ctx, getJob, id)
	var i Job
	err := row.Scan(&i.ID, &i.Name, &i.Cron)
	return i, err
}

const getJobsView = `-- name: GetJobsView :many
SELECT
    id, name, cron, runs
FROM
    jobs_view
ORDER BY
    cron,
    name
`

func (q *Queries) GetJobsView(ctx context.Context) ([]JobsView, error) {
	rows, err := q.db.QueryContext(ctx, getJobsView)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []JobsView
	for rows.Next() {
		var i JobsView
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Cron,
			&i.Runs,
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

const listJobs = `-- name: ListJobs :many
SELECT
    id, name, cron
FROM
    jobs
ORDER BY
    cron,
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
		if err := rows.Scan(&i.ID, &i.Name, &i.Cron); err != nil {
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
SET
    name = ?,
    cron = ?
WHERE
    id = ?
`

type UpdateJobParams struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
	ID   string `json:"id"`
}

func (q *Queries) UpdateJob(ctx context.Context, arg UpdateJobParams) error {
	_, err := q.db.ExecContext(ctx, updateJob, arg.Name, arg.Cron, arg.ID)
	return err
}
