-- name: CreateSession :exec
INSERT INTO sessions (jti, email, username, created_at, expires_at, revoked)
VALUES (?, ?, ?, ?, ?, 0);

-- name: GetActiveSession :one
SELECT
  jti,
  email,
  username,
  created_at,
  expires_at,
  revoked
FROM
  sessions
WHERE
  jti = ?
  AND revoked = 0
  AND expires_at > ?;

-- name: RevokeSession :exec
UPDATE sessions
SET
  revoked = 1
WHERE
  jti = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE
  expires_at < ?;
