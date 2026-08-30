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
