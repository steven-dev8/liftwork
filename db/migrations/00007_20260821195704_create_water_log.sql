-- +goose Up
CREATE TABLE water_logs (
    id BIGSERIAL PRIMARY KEY,

    water_goal_id BIGINT NOT NULL
        REFERENCES water_goal(id)
        ON DELETE CASCADE,

    amount_ml INTEGER NOT NULL
        CHECK (amount_ml > 0),

    consumed_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_water_logs_water_goal_id
ON water_logs(water_goal_id);


-- +goose Down
DROP TABLE IF EXISTS water_logs;
