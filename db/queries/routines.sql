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

-- name: ListRoutine :many
SELECT *
FROM routines
WHERE user_id = @user_id;

-- name: GetExerciseRoutine :many
SELECT
    e.name,
    re.position,
    re.target_sets,
    re.target_reps_min,
    re.target_reps_max
FROM routine_exercises re
INNER JOIN exercises e
    ON e.id = re.exercise_id
WHERE re.routine_id = @routine_id
ORDER BY re.position;
