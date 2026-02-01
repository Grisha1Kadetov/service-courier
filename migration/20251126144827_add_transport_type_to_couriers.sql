-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS couriers
    ADD COLUMN IF NOT EXISTS transport_type TEXT NOT NULL DEFAULT 'on_foot';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS couriers
    DROP COLUMN IF EXISTS transport_type;
-- +goose StatementEnd
