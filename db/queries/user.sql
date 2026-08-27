-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash
)
VALUES (
    sqlc.narg(email),
    @username,
    @password_hash
)
RETURNING id, username, created_at;


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


-- name: GetUser :one
SELECT id, username, password_hash
FROM users
WHERE username = @username;
