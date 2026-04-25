-- +goose Up
ALTER TABLE price_restore_batches
    ADD COLUMN IF NOT EXISTS shop_id            BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS restore_status     TEXT   NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS restore_started_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_price_restore_batches_shop_id
    ON price_restore_batches (shop_id);

CREATE INDEX IF NOT EXISTS idx_price_restore_batches_restore_status
    ON price_restore_batches (restore_status);

-- +goose Down
DROP INDEX IF EXISTS idx_price_restore_batches_restore_status;
DROP INDEX IF EXISTS idx_price_restore_batches_shop_id;

ALTER TABLE price_restore_batches
    DROP COLUMN IF EXISTS restore_started_at,
    DROP COLUMN IF EXISTS restore_status,
    DROP COLUMN IF EXISTS shop_id;
