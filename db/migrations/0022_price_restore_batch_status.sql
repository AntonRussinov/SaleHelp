-- +goose Up
ALTER TABLE price_restore_batches
    ADD COLUMN IF NOT EXISTS shop_id         BIGINT,
    ADD COLUMN IF NOT EXISTS restore_status  TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS restore_started_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE price_restore_batches
    DROP COLUMN IF EXISTS restore_started_at,
    DROP COLUMN IF EXISTS restore_status,
    DROP COLUMN IF EXISTS shop_id;
