-- name: ListRunsByJobID :many
SELECT
    *
FROM
    runs
WHERE
    job_id = ?
ORDER BY
    id DESC;

-- name: GetRunsView :many
SELECT
    *
FROM
    runs_view
WHERE
    job_id = ?
ORDER BY
    start_time DESC
LIMIT
    ?;

-- name: CreateRun :one
INSERT INTO
    runs (job_id, status_id)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateRun :exec
UPDATE runs
SET
    status_id = ?,
    end_time = ?
WHERE
    id = ?;