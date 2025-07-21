-- +goose Up
-- +goose StatementBegin
CREATE TYPE link_status_enum AS ENUM ('available', 'full', 'closed', 'invalid');
-- +goose StatementEnd