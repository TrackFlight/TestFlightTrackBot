-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_links (
    chat_id              BIGINT NOT NULL CONSTRAINT fk_chat_links_chat REFERENCES chats ON DELETE CASCADE,
    link_id              BIGINT NOT NULL CONSTRAINT fk_chat_links_link REFERENCES links ON DELETE CASCADE,
    allow_opened         BOOLEAN DEFAULT true NOT NULL,
    allow_closed         BOOLEAN DEFAULT false NOT NULL,
    last_notified_status link_status_enum,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, link_id)
);

CREATE INDEX IF NOT EXISTS idx_chat_links_updated_at
    ON chat_links (updated_at);

CREATE INDEX IF NOT EXISTS idx_chat_links_created_at
    ON chat_links (created_at);

CREATE INDEX IF NOT EXISTS idx_chat_links_last_notified_status
    ON chat_links (last_notified_status);

CREATE INDEX IF NOT EXISTS idx_chat_links_allow_closed
    ON chat_links (allow_closed);

CREATE INDEX IF NOT EXISTS idx_chat_links_allow_opened
    ON chat_links (allow_opened);
-- +goose StatementEnd