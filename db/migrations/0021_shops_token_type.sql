-- +goose Up
ALTER TABLE shops
    ADD COLUMN IF NOT EXISTS token_type TEXT NOT NULL DEFAULT 'basic';

-- +goose Down
ALTER TABLE shops DROP COLUMN IF EXISTS token_type;
