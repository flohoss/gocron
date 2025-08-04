-- name: ListLogsByRunIDs :many
SELECT
    created_at,
    run_id,
    severity_id,
    message,
    STRFTIME(
        '%Y-%m-%d %H:%M:%S',
        created_at / 1000,
        'unixepoch',
        'localtime'
    ) AS created_at_time
FROM
    logs
WHERE
    run_id IN (sqlc.slice (run_ids))
ORDER BY
    run_id,
    created_at;

-- name: CreateLog :one
INSERT INTO
    logs (created_at, run_id, severity_id, message)
VALUES
    (?, ?, ?, ?) RETURNING *;