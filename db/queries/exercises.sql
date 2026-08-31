-- name: CreateExercise :one
INSERT INTO exercises(
    user_id,
    name,
    muscle_group,
    notes
) VALUES (
    @user_id,
    @name,
    @muscle_group,
    @notes
)
RETURNING id, user_id, name, muscle_group, notes, created_at, updated_at;

-- name: GetExercises :many
SELECT id, name, muscle_group, notes, created_at, updated_at
FROM exercises
WHERE user_id = @user_id;

-- name: UpdateExerciseById :one
UPDATE exercises
SET
    name = COALESCE(sqlc.narg(name), name),
    muscle_group = COALESCE(sqlc.narg(muscle_group), muscle_group),
    notes = COALESCE(sqlc.narg(notes), notes),
    updated_at = now()
WHERE id = @id AND user_id = @user_id
RETURNING *;