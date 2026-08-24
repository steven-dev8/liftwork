-- +goose Up
CREATE TABLE routines (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(1) NOT NULL UNIQUE CHECK (code IN ('A', 'B', 'C', 'D')),
    name VARCHAR(120) NOT NULL CHECK (btrim(name) <> ''),
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS routines;