-- +goose Up
CREATE UNIQUE INDEX one_open_workout_per_user
ON workout_sessions (user_id)
WHERE finished_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS one_open_workout_per_user;
