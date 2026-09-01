-- name: CreateRoutine :one
INSERT INTO routines (
    user_id,
    code,
    name,
    description,
    created_at,
    updated_at
)
VALUES (
    @user_id,
    @code,
    @name,
    @description,
    now(),
    now()
)
RETURNING *;

