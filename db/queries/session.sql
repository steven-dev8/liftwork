-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    refresh_token_hash,
    expires_at
) VALUES (
    @user_id,
    @refresh_token_hash,
    @expires_at
)
RETURNING user_id, refresh_token_hash, created_at, expires_at;

-- name: GetSessionByRefreshTokenHash :one
SELECT id, user_id, revoked_at, expires_at
FROM sessions
WHERE refresh_token_hash = @refresh_token_hash;

-- name: RevokeSession :execrows
UPDATE sessions
SET revoked_at = now()
WHERE refresh_token_hash = @refresh_token_hash
  AND revoked_at IS NULL;

-- name: RotateRefreshToken :execrows
UPDATE sessions
SET refresh_token_hash = @new_refresh_token_hash
WHERE id = @session_id
  AND refresh_token_hash = @old_refresh_token_hash
  AND revoked_at IS NULL
  AND expires_at > now();
