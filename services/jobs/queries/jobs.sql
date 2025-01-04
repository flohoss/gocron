-- name: ListJobs :many
SELECT
    *
FROM
    jobs
ORDER BY
    name;

-- name: ListJobsAndLatestRun :many
SELECT
    j.id,
    j.name,
    j.cron,
    r.start_time,
    r.end_time,
    r.status_id
FROM
    jobs j
    LEFT JOIN runs r ON j.id = r.job_id
    AND r.id = (
        SELECT
            MAX(id)
        FROM
            runs
        WHERE
            runs.job_id = j.id
    )
ORDER BY
    name;

-- name: GetJob :one
SELECT
    *
FROM
    jobs
WHERE
    id = ?;

-- name: CreateJob :one
INSERT INTO
    jobs (id, name, cron)
VALUES
    (?, ?, ?) RETURNING *;

-- name: UpdateJob :exec
UPDATE jobs
SET
    name = ?,
    cron = ?
WHERE
    id = ?;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE
    id = ?;