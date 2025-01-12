-- name: GetRunsView4 :many
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
            4
    ) subquery
ORDER BY
    start_time ASC;

-- name: GetRunsView20 :many
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