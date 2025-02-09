-- name: GetRunsViewHome :many
SELECT
    *
FROM
    (
        SELECT
            *
        FROM
            runs_view
        WHERE
            job_id = ?
        ORDER BY
            start_time DESC
        LIMIT
            3
    ) subquery
ORDER BY
    start_time ASC;

-- name: GetRunsViewDetail :many
SELECT
    *
FROM
    runs_view
WHERE
    job_id = ?
ORDER BY
    start_time DESC
LIMIT
    20;

-- name: CreateRun :one
INSERT INTO
    runs (job_id, status_id, start_time)
VALUES
    (?, ?, ?) RETURNING *;

-- name: UpdateRun :exec
UPDATE runs
SET
    status_id = ?,
    end_time = ?
WHERE
    id = ?;

-- name: IsIdle :one
SELECT
    CASE
        WHEN status_id = 1 THEN FALSE
        ELSE TRUE
    END AS is_idle
FROM
    runs
UNION ALL
SELECT
    TRUE
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            runs
    )
LIMIT
    1;