-- name: GetSchemaVersion :one
SELECT
  version
FROM
  schema_version
LIMIT
  1;

-- name: SetSchemaVersion :exec
INSERT
OR IGNORE INTO schema_version (version)
VALUES
  (?);
