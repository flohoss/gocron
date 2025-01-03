-- name: ListStatus :many
SELECT
    *
FROM
    status
ORDER BY
    id;

-- name: CreateStatus :one
INSERT INTO
    status (status)
VALUES
    (?) RETURNING *