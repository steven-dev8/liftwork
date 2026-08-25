-- +goose Up
CREATE TABLE performed_sets (
    id BIGSERIAL PRIMARY KEY,

    workout_exercise_id BIGINT NOT NULL
        REFERENCES workout_exercises(id)
        ON DELETE CASCADE,

    set_number INTEGER NOT NULL
        CHECK (set_number > 0),

    weight_kg NUMERIC(7,2) NOT NULL
        CHECK (weight_kg >= 0),

    repetitions INTEGER NOT NULL
        CHECK (repetitions > 0),

    rir INTEGER NOT NULL
        CHECK (rir BETWEEN 0 AND 10),

    performed_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (workout_exercise_id, set_number)
);

-- +goose Down
DROP TABLE IF EXISTS performed_sets;
