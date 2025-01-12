-- name: ListJobs :many
SELECT
    *
FROM
    jobs
ORDER BY
    cron,
    name;

-- name: GetJobsView :many
SELECT
    *
FROM
    jobs_view
ORDER BY
    cron,
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
    jobs (id, name, cron)
VALUES
    (?, ?, ?) RETURNING *;

-- name: UpdateJob :exec
UPDATE jobs
SET
    name = ?,
    cron = ?
WHERE
    id = ?;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    id = ?;