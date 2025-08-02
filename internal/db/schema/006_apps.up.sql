-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS apps (
    id          BIGSERIAL PRIMARY KEY,
    app_name    VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    icon_url    VARCHAR(255),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_apps_updated_at
    ON apps(updated_at);

CREATE INDEX IF NOT EXISTS idx_apps_created_at
    ON apps(created_at);

CREATE INDEX IF NOT EXISTS idx_apps_app_name
    ON apps(app_name);
-- +goose StatementEnd