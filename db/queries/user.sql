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

-- name: GetUser :one
SELECT id, username, password_hash
FROM users
WHERE username = @username;
