-- name: GetRetentionPolicies :one
SELECT
    *
FROM
    retention_policies
WHERE
    retention_policy_id = ?
LIMIT
    1;

-- name: ListRetentionPolicies :many
SELECT
    *
FROM
    retention_policies
ORDER BY
    retention_policy;

-- name: CreateRetentionPolicy :one
INSERT INTO
    retention_policies (retention_policy)
VALUES
    (?) RETURNING *;

-- name: UpdateRetentionPolicy :one
UPDATE retention_policies
SET
    retention_policy = ?
WHERE
    retention_policy_id = ? RETURNING *;

-- name: DeleteRetentionPolicy :exec
DELETE FROM retention_policies
WHERE
    retention_policy_id = ?;