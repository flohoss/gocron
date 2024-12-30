-- name: ListJobs :many
SELECT
    *
FROM
    jobs
ORDER BY
    name;

-- name: GetJob :one
SELECT
    *
FROM
    jobs
WHERE
    id = ?;

-- name: CreateJob :one
INSERT INTO
    jobs (Name, backup_schedule_id)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateJob :exec
UPDATE jobs
set
    name = ?,
    backup_schedule_id = ?
WHERE
    id = ?;