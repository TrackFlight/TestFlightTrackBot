-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS preferences (
    chat_id                         BIGINT PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
    new_features_notifications      BOOLEAN DEFAULT true NOT NULL,
    weekly_insights_notifications   BOOLEAN DEFAULT true NOT NULL,
    created_at                      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_preferences_weekly_insights_notifications
    ON preferences (weekly_insights_notifications);

CREATE INDEX IF NOT EXISTS idx_preferences_new_features_notifications
    ON preferences (new_features_notifications);
-- +goose StatementEnd