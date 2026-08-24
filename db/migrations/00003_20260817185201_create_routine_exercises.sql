-- +goose Up
CREATE TABLE routine_exercises (
    routine_id BIGINT NOT NULL REFERENCES routines (id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises (id) ON DELETE RESTRICT,
    position INTEGER NOT NULL CHECK (position > 0),
    target_sets INTEGER NOT NULL DEFAULT 3 CHECK (target_sets > 0),
    target_reps_min INTEGER NOT NULL DEFAULT 1 CHECK (target_reps_min > 0),
    target_reps_max INTEGER NOT NULL DEFAULT 1 CHECK (target_reps_max > 0),
    PRIMARY KEY (routine_id, exercise_id),
    UNIQUE (routine_id, position)
);

-- +goose Down
DROP TABLE IF EXISTS routine_exercises;
