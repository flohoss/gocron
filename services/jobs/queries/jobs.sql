-- name: ListJobs :many
SELECT
    *
FROM
    jobs
ORDER BY
    name;

-- name: ListJobsCommandsEnvsRunsAndLogs :many
SELECT
    sqlc.embed(jobs),
    sqlc.embed(commands),
    sqlc.embed(envs),
    sqlc.embed(runs)
FROM
    jobs
    JOIN commands ON jobs.id = commands.job_id
    JOIN envs ON jobs.id = envs.job_id
    JOIN runs ON jobs.id = runs.job_id
WHERE
    jobs.id = ?
ORDER BY
    jobs.name;

-- name: ListJobsWithLatestRun :many
WITH
    latest_runs AS (
        SELECT
            job_id,
            MAX(id) AS max_run_id
        FROM
            runs
        GROUP BY
            job_id
    )
SELECT
    sqlc.embed(jobs),
    sqlc.embed(runs)
FROM
    jobs
    JOIN latest_runs lr ON jobs.id = lr.job_id
    JOIN runs ON lr.max_run_id = runs.id
ORDER BY
    jobs.name;

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