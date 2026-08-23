-- +goose Up
CREATE TABLE exercises (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT
        REFERENCES users(id)
        ON DELETE CASCADE,
    name VARCHAR(160) NOT NULL CHECK (btrim(name) <> ''),
    muscle_group VARCHAR(120) NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE NULLS NOT DISTINCT (user_id, name)
);

-- +goose Down
DROP TABLE IF EXISTS exercises;
