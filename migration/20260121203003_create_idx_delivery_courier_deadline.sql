-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_delivery_courier_deadline 
ON delivery (courier_id, deadline);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_delivery_courier_deadline;
