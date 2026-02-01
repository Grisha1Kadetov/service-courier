-- +goose NO TRANSACTION

-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_couriers_status ON couriers (status);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_couriers_status;
