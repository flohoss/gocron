-- name: ListSeverities :many
SELECT
    *
FROM
    severities
ORDER BY
    id;

-- name: CreateSeverity :one
INSERT INTO
    severities (severity)
VALUES
    (?) RETURNING *;