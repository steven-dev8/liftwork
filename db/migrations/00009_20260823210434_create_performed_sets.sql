-- +goose Up
CREATE TABLE performed_sets (
    id BIGSERIAL PRIMARY KEY,
    workout_session_id BIGINT NOT NULL REFERENCES workout_sessions (id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises (id) ON DELETE RESTRICT,
    set_number INTEGER NOT NULL CHECK (set_number > 0),
    weight_kg NUMERIC(7, 2) NOT NULL CHECK (weight_kg >= 0),
    repetitions INTEGER NOT NULL CHECK (repetitions > 0),
    rir INTEGER NOT NULL CHECK (rir BETWEEN 0 AND 10),
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workout_session_id, exercise_id, set_number)
);

-- +goose Down
DROP TABLE IF EXISTS performed_sets;
