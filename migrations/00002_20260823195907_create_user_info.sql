-- +goose Up
CREATE TABLE user_info(
    user_id BIGINT PRIMARY KEY
        REFERENCES users(id)
        ON DELETE CASCADE,
    weight NUMERIC(6, 2),
    avatar TEXT,
    date_of_birth DATE
);

-- +goose Down
DROP TABLE IF EXISTS user_info;
