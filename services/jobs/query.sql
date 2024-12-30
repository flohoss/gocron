-- name: ListJobs :many
SELECT
    *
FROM
    jobs
ORDER BY
    name;

-- name: GetJob :one
SELECT
    *
FROM
    jobs
WHERE
    name = ?;

-- name: CreateJob :one
INSERT INTO
    jobs (name)
VALUES
    (?) RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    name = ?;