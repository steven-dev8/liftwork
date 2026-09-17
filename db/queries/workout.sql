-- name: CreateWorkoutSessionWithoutRoutine :one
INSERT INTO workout_sessions (
    user_id,
    routine_id,
    notes,
    created_at
)
VALUES (
    @user_id,
    NULL,
    @notes,
    now()
)
RETURNING *;


-- name: CreateWorkoutSessionWithRoutine :one
INSERT INTO workout_sessions (
    user_id,
    routine_id,
    notes,
    started_at,
    created_at
)
SELECT
    @user_id,
    r.id,
    @notes,
    now(),
    now()
FROM routines r
WHERE r.id = @routine_id
  AND r.user_id = @user_id
RETURNING *;


-- name: CreateWorkoutExercisesFromRoutine :many
WITH inserted AS (
    INSERT INTO workout_exercises (
        workout_session_id,
        exercise_id,
        position
    )
    SELECT
        @workout_session_id,
        re.exercise_id,
        re.position
    FROM routine_exercises re
    WHERE re.routine_id = @routine_id
    ORDER BY re.position
    RETURNING
        id,
        workout_session_id,
        exercise_id,
        position
)
SELECT
    i.id,
    i.exercise_id,
    i.position,
    e.name,
    e.muscle_group,
    e.notes,
    re.target_sets,
    re.target_reps_min,
    re.target_reps_max
FROM inserted i
JOIN exercises e
    ON e.id = i.exercise_id
JOIN routine_exercises re
    ON re.routine_id = @routine_id
    AND re.exercise_id = i.exercise_id
ORDER BY i.position;
