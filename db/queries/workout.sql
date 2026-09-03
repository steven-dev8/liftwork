-- name: CreateWorkoutSession :one
INSERT INTO workout_sessions (
    user_id,
    routine_id,
    notes,
    created_at
)
SELECT
    @user_id,
    r.id,
    @notes,
    now()
FROM routines r
WHERE r.id = @routine_id
  AND r.user_id = @user_id
RETURNING *;
