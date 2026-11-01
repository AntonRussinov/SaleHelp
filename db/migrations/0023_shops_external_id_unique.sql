-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS uidx_shops_external_id
    ON shops (external_id)
    WHERE external_id != '';

-- +goose Down
DROP INDEX IF EXISTS uidx_shops_external_id;
