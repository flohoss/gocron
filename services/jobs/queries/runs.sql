-- name: ListRuns :many
SELECT
    *
FROM
    runs
WHERE
    job_id = ?
ORDER BY
    start_time DESC;

-- name: ListRunsAndLogs :many
SELECT
    *
FROM
    runs
    LEFT JOIN logs ON runs.id = logs.run_id
WHERE
    job_id = ?
ORDER BY
    start_time DESC;

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