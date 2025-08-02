-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id                BIGSERIAL PRIMARY KEY,
    url               VARCHAR(255) NOT NULL UNIQUE,
    app_id            BIGINT REFERENCES apps(id) ON DELETE CASCADE,
    status            link_status_enum,
    is_public         BOOLEAN NOT NULL DEFAULT FALSE,
    last_availability TIMESTAMP WITH TIME ZONE,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_links_updated_at
    ON links (updated_at);

CREATE INDEX IF NOT EXISTS idx_links_created_at
    ON links (created_at);

CREATE INDEX IF NOT EXISTS idx_links_status
    ON links (status);

CREATE INDEX IF NOT EXISTS idx_links_app_id
    ON links (app_id);

CREATE INDEX IF NOT EXISTS idx_links_url
    ON links (url);
-- +goose StatementEnd