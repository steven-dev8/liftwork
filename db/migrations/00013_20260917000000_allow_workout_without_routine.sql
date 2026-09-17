-- +goose Up
ALTER TABLE workout_sessions
DROP CONSTRAINT workout_sessions_check1;

-- +goose Down
ALTER TABLE workout_sessions
ADD CONSTRAINT workout_sessions_check1
CHECK (routine_id IS NOT NULL OR finished_at IS NOT NULL);
