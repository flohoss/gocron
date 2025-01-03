-- name: ListLogs :many
SELECT
    *
FROM
    logs
WHERE
    run_id = ?
ORDER BY
    created_at DESC;

-- name: CreateLog :one
INSERT INTO
    logs (run_id, severity_id, message)
VALUES
    (?, ?, ?) RETURNING *;