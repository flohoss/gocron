-- name: ListEnvs :many
SELECT
    *
FROM
    envs
ORDER BY
    id;

-- name: ListEnvsByJobID :many
SELECT
    *
FROM
    envs
WHERE
    job_id = ?
ORDER BY
    id;

-- name: CreateEnv :one
INSERT INTO
    envs (id, job_id, KEY, value)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: DeleteEnvs :exec
DELETE FROM envs;