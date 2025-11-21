-- +goose Up
-- +goose StatementBegin
CREATE TYPE link_status_enum AS ENUM ('available', 'full', 'closed', 'invalid');

CREATE TYPE user_status_enum AS ENUM ('reachable', 'blocked_by_user', 'deleted_account');
-- +goose StatementEnd