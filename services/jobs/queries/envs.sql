-- name: ListEnvs :many
SELECT
    *
FROM
    envs
ORDER BY
    job_id,
    key;

-- name: ListEnvsByJobID :many
SELECT
    *
FROM
    envs
WHERE
    job_id = ?
ORDER BY
    key;

-- name: CreateEnv :one
INSERT INTO
    envs (job_id, KEY, value)
VALUES
    (?, ?, ?) RETURNING *;

-- name: DeleteEnvs :exec
DELETE FROM envs;