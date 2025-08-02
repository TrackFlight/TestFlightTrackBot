-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS translator_languages (
    chat_id    BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    lang       VARCHAR(5) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, lang)
);
-- +goose StatementEnd