-- name: GetRuns :many
SELECT
    id,
    job_name,
    status_id,
    start_time,
    end_time
FROM
    runs
WHERE
    job_name_normalized = ?
ORDER BY
    start_time DESC
LIMIT
    ?;

-- name: GetThreeRunsPerJobName :many
WITH
    ranked_runs AS (
        SELECT
            id,
            job_name,
            job_name_normalized,
            status_id,
            start_time,
            end_time,
            ROW_NUMBER() OVER (
                PARTITION BY
                    job_name_normalized
                ORDER BY
                    start_time DESC
            ) AS rn
        FROM
            runs
    )
SELECT
    id,
    job_name,
    status_id,
    start_time,
    end_time
FROM
    ranked_runs
WHERE
    rn <= 3
ORDER BY
    job_name_normalized,
    start_time DESC;

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
    job_name_normalized NOT IN (sqlc.slice (job_names));