-- name: CreateSession :exec
INSERT INTO sessions (jti, subject, created_at, expires_at)
VALUES (?, ?, ?, ?);

-- name: GetActiveSession :one
SELECT
  jti,
  subject,
  created_at,
  expires_at
FROM
  sessions
WHERE
  jti = ?
  AND expires_at > ?;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
  jti = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE
  expires_at < ?;
