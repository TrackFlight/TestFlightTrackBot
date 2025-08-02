-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS premium_users (
    chat_id    BIGINT PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_premium_users_updated_at
    ON premium_users (updated_at);

CREATE INDEX IF NOT EXISTS idx_premium_users_created_at
    ON premium_users (created_at);
-- +goose StatementEnd