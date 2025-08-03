-- name: GetRuns :many
SELECT
    *
FROM
    runs
WHERE
    job_name = ?
ORDER BY
    start_time DESC
LIMIT
    ?;

-- name: CreateRun :one
INSERT INTO
    runs (job_name, status_id, start_time)
VALUES
    (?, ?, ?) RETURNING *;

-- name: UpdateRun :one
UPDATE runs
SET
    status_id = ?,
    end_time = ?
WHERE
    id = ? RETURNING *;

-- name: IsIdle :one
SELECT
    CAST(
        NOT EXISTS (
            SELECT
                1
            FROM
                runs
            WHERE
                status_id = 1
        ) AS INTEGER
    ) AS is_idle;

-- name: DeleteOldRuns :exec
DELETE FROM runs
WHERE
    start_time < ?;

-- name: DeleteObsoleteRuns :exec
DELETE FROM runs
WHERE
    job_name NOT IN (sqlc.slice (job_names));