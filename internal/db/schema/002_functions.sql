-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;

CREATE OR REPLACE FUNCTION assert(
    cond boolean,
    msg text,
    err_code text DEFAULT 'P0001'
)
    RETURNS void AS $$
BEGIN
    IF NOT cond THEN
        RAISE EXCEPTION '%', msg
            USING ERRCODE = err_code;
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd