-- +goose Up
CREATE TABLE workout_exercises (
    id BIGSERIAL PRIMARY KEY,

    workout_session_id BIGINT NOT NULL
        REFERENCES workout_sessions(id)
        ON DELETE CASCADE,

    exercise_id BIGINT NOT NULL
        REFERENCES exercises(id)
        ON DELETE RESTRICT,

    position INTEGER NOT NULL
        CHECK (position > 0),

    UNIQUE (workout_session_id, position)
);

-- +goose Down
DROP TABLE IF EXISTS workout_exercises;
