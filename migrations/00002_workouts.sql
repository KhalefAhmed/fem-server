-- +goose UP
-- +goose statementBegin
CREATE TABLE IF NOT EXISTS workouts
(
    id               BIGINT PRIMARY KEY,
    -- user id
    description      TEXT,
    duration_minutes INTEGER NOT NULL,
    calories_burned  INTEGER,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose statementEnd

-- +goose DOWN
-- +goose statementBegin
DROP TABLE workouts;
-- +goose statementEnd
