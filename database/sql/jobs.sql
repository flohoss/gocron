-- name: GetJob :one
SELECT
    *
FROM
    jobs
WHERE
    job_id = ?
LIMIT
    1;

-- name: ListJobs :many
SELECT
    *
FROM
    jobs
ORDER BY
    description;

-- name: CreateJob :one
INSERT INTO
    jobs (
        description,
        local_directory,
        restic_remote,
        restart_option,
        password_file_path,
        compression_type,
        svg_icon,
        retention_policy_id
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateJob :one
UPDATE jobs
SET
    description = ?,
    local_directory = ?,
    restic_remote = ?,
    restart_option = ?,
    password_file_path = ?,
    compression_type = ?,
    svg_icon = ?,
    retention_policy_id = ?
WHERE
    job_id = ? RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    job_id = ?;