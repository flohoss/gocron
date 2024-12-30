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
    name = ?;

-- name: CreateJob :one
INSERT INTO
    jobs (name, cron)
VALUES
    (?, ?) ON CONFLICT (name) DO
UPDATE
SET
    cron = EXCLUDED.cron RETURNING *;

-- name: UpdateJob :exec
UPDATE jobs
SET
    cron = ?
WHERE
    name = ?;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    name = ?;