-- +goose Up
CREATE TABLE routines (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL CHECK (code IN ('A', 'B', 'C', 'D', 'E')),
    name VARCHAR(120) NOT NULL CHECK (btrim(name) <> ''),
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (user_id, code)
);

-- +goose Down
DROP TABLE IF EXISTS routines;
