// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: jobs.sql

package database

import (
	"context"
)

const createJob = `-- name: CreateJob :one
INSERT INTO
    jobs (
        description,
        local_directory,
        restic_remote,
        restart_option,
        password_file_path,
        compression_type,
        svg_icon,
        retention_policy_id
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?) RETURNING job_id, description, local_directory, restic_remote, restart_option, password_file_path, compression_type, svg_icon, created_at, retention_policy_id
`

type CreateJobParams struct {
	Description       string
	LocalDirectory    string
	ResticRemote      string
	RestartOption     int64
	PasswordFilePath  string
	CompressionType   string
	SvgIcon           string
	RetentionPolicyID int64
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, createJob,
		arg.Description,
		arg.LocalDirectory,
		arg.ResticRemote,
		arg.RestartOption,
		arg.PasswordFilePath,
		arg.CompressionType,
		arg.SvgIcon,
		arg.RetentionPolicyID,
	)
	var i Job
	err := row.Scan(
		&i.JobID,
		&i.Description,
		&i.LocalDirectory,
		&i.ResticRemote,
		&i.RestartOption,
		&i.PasswordFilePath,
		&i.CompressionType,
		&i.SvgIcon,
		&i.CreatedAt,
		&i.RetentionPolicyID,
	)
	return i, err
}

const deleteJob = `-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    job_id = ?
`

func (q *Queries) DeleteJob(ctx context.Context, jobID int64) error {
	_, err := q.db.ExecContext(ctx, deleteJob, jobID)
	return err
}

const getJob = `-- name: GetJob :one
SELECT
    job_id, description, local_directory, restic_remote, restart_option, password_file_path, compression_type, svg_icon, created_at, retention_policy_id
FROM
    jobs
WHERE
    job_id = ?
LIMIT
    1
`

func (q *Queries) GetJob(ctx context.Context, jobID int64) (Job, error) {
	row := q.db.QueryRowContext(ctx, getJob, jobID)
	var i Job
	err := row.Scan(
		&i.JobID,
		&i.Description,
		&i.LocalDirectory,
		&i.ResticRemote,
		&i.RestartOption,
		&i.PasswordFilePath,
		&i.CompressionType,
		&i.SvgIcon,
		&i.CreatedAt,
		&i.RetentionPolicyID,
	)
	return i, err
}

const listJobs = `-- name: ListJobs :many
SELECT
    job_id, description, local_directory, restic_remote, restart_option, password_file_path, compression_type, svg_icon, created_at, retention_policy_id
FROM
    jobs
ORDER BY
    description
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
		if err := rows.Scan(
			&i.JobID,
			&i.Description,
			&i.LocalDirectory,
			&i.ResticRemote,
			&i.RestartOption,
			&i.PasswordFilePath,
			&i.CompressionType,
			&i.SvgIcon,
			&i.CreatedAt,
			&i.RetentionPolicyID,
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

const updateJob = `-- name: UpdateJob :one
UPDATE jobs
SET
    description = ?,
    local_directory = ?,
    restic_remote = ?,
    restart_option = ?,
    password_file_path = ?,
    compression_type = ?,
    svg_icon = ?,
    retention_policy_id = ?
WHERE
    job_id = ? RETURNING job_id, description, local_directory, restic_remote, restart_option, password_file_path, compression_type, svg_icon, created_at, retention_policy_id
`

type UpdateJobParams struct {
	Description       string
	LocalDirectory    string
	ResticRemote      string
	RestartOption     int64
	PasswordFilePath  string
	CompressionType   string
	SvgIcon           string
	RetentionPolicyID int64
	JobID             int64
}

func (q *Queries) UpdateJob(ctx context.Context, arg UpdateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, updateJob,
		arg.Description,
		arg.LocalDirectory,
		arg.ResticRemote,
		arg.RestartOption,
		arg.PasswordFilePath,
		arg.CompressionType,
		arg.SvgIcon,
		arg.RetentionPolicyID,
		arg.JobID,
	)
	var i Job
	err := row.Scan(
		&i.JobID,
		&i.Description,
		&i.LocalDirectory,
		&i.ResticRemote,
		&i.RestartOption,
		&i.PasswordFilePath,
		&i.CompressionType,
		&i.SvgIcon,
		&i.CreatedAt,
		&i.RetentionPolicyID,
	)
	return i, err
}
