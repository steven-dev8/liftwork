-- +goose Up
CREATE TABLE water_goal (
    user_id BIGINT PRIMARY KEY
        REFERENCES users(id)
        ON DELETE CASCADE,
    goal_amount_ml INTEGER NOT NULL
        CHECK (goal_amount_ml > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS water_goal;



