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
    jobs (name)
VALUES
    (?) RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    name = ?;

-- name: CreateRun :one
INSERT INTO
    runs (job, status_id)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateRun :exec
UPDATE runs
SET
    status_id = ?,
    end_time = ?
WHERE
    id = ?;

-- name: CreateLog :one
INSERT INTO
    logs (run_id, severity_id, message)
VALUES
    (?, ?, ?) RETURNING *;

-- name: ListSeverities :many
SELECT
    *
FROM
    severities
ORDER BY
    id;

-- name: CreateSeverity :one
INSERT INTO
    severities (severity)
VALUES
    (?) RETURNING *;

-- name: ListStatus :many
SELECT
    *
FROM
    status
ORDER BY
    id;

-- name: CreateStatus :one
INSERT INTO
    status (status)
VALUES
    (?) RETURNING *;