-- name: ListLogsByRunID :many
SELECT
    *,
    STRFTIME(
        '%Y-%m-%d %H:%M:%S',
        created_at / 1000,
        'unixepoch',
        'localtime'
    ) AS created_at_time
FROM
    logs
WHERE
    run_id = ?
ORDER BY
    created_at;

-- name: CreateLog :one
INSERT INTO
    logs (created_at, run_id, severity_id, message)
VALUES
    (?, ?, ?, ?) RETURNING *;