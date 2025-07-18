-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    CREATE TYPE link_status_enum AS ENUM ('available', 'full', 'closed', 'invalid');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +goose StatementEnd