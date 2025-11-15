-- +goose Up
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_created_at ON questions(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_questions_created_at;
DROP TABLE IF EXISTS questions CASCADE;

