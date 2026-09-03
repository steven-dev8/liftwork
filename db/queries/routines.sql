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
    e.id,
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


-- name: AddExerciseRoutine :execrows
INSERT INTO routine_exercises (
    routine_id,
    exercise_id,
    position,
    target_sets,
    target_reps_min,
    target_reps_max
)
SELECT
    r.id,
    e.id,
    @position,
    @target_sets,
    @target_reps_min,
    @target_reps_max
FROM routines r
JOIN exercises e
    ON e.id = @exercise_id
WHERE r.id = @routine_id
  AND r.user_id = @user_id
  AND (
      e.user_id = @user_id
      OR e.user_id IS NULL
  );

-- name: DeleteExerciseRoutine :one
DELETE FROM routine_exercises re
USING routines r
WHERE re.routine_id = @routine_id
  AND re.exercise_id = @exercise_id
  AND r.id = re.routine_id
  AND r.user_id = @user_id
RETURNING re.position;

-- name: ReorderRoutineExercises :exec
UPDATE routine_exercises
SET position = position - 1
WHERE routine_id = @routine_id
  AND position > @deleted_position;

-- name: DeleteRoutine :execrows
DELETE FROM routines
WHERE id = @id AND user_id = @user_id;

-- name: UpdateRoutine :one
UPDATE routines
SET
    code = COALESCE(sqlc.narg(code), code),
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at = now()
WHERE id = @id
  AND user_id = @user_id
RETURNING *;
