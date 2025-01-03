-- name: ListCommands :many
SELECT
    *
FROM
    commands
ORDER BY
    job_id;

-- name: ListCommandsByJobID :many
SELECT
    *
FROM
    commands
WHERE
    job_id = ?;

-- name: CreateCommand :one
INSERT INTO
    commands (job_id, command)
VALUES
    (?, ?) RETURNING *;

-- name: DeleteCommands :exec
DELETE FROM commands;