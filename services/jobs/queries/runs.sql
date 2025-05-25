-- name: GetRunsView :many
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
            ?
    ) subquery
ORDER BY
    start_time ASC;

-- name: CreateRun :one
INSERT INTO
    runs (job_id, status_id, start_time)
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

-- name: DeleteRuns :exec
DELETE FROM runs
WHERE
    start_time < ?;