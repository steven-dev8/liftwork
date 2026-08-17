-- +goose Up
CREATE TABLE water_logs (
    id BIGSERIAL PRIMARY KEY,
    amount_ml INTEGER NOT NULL CHECK (amount_ml > 0),
    consumed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS water_logs;
