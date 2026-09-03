-- +goose Up
CREATE TABLE workout_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,
    routine_id BIGINT
        REFERENCES routines(id)
        ON DELETE SET NULL,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (finished_at IS NULL OR finished_at >= started_at),
    CHECK (routine_id IS NOT NULL OR finished_at IS NOT NULL)
);

-- +goose Down
DROP TABLE IF EXISTS workout_sessions;
