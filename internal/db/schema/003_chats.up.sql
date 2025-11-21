-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chats (
    id         BIGINT PRIMARY KEY,
    lang       VARCHAR(10) DEFAULT 'en'::CHARACTER VARYING NOT NULL,
    status     user_status_enum NOT NULL DEFAULT 'reachable'::user_status_enum,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chats_status
    ON chats (status);

CREATE INDEX IF NOT EXISTS idx_chats_updated_at
    ON chats (updated_at);

CREATE INDEX IF NOT EXISTS idx_chats_created_at
    ON chats (created_at);
-- +goose StatementEnd