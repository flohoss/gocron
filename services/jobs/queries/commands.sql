-- name: ListCommands :many
SELECT
    *
FROM
    commands
ORDER BY
    id;

-- name: ListCommandsByJobID :many
SELECT
    *
FROM
    commands
WHERE
    job_id = ?;

-- name: CreateCommand :one
INSERT INTO
    commands (id, job_id, command, file_output)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: DeleteCommands :exec
DELETE FROM commands;