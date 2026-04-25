-- +goose Up
-- +goose StatementBegin

-- Поля для разделения синка на быстрый (каждые 30 мин) и тяжёлый (раз в час).
-- Тяжёлый синк: карточки товаров, цены, остатки, склады.
-- Инкрементальный синк карточек через cursor WB API (updatedAt + nmID).

ALTER TABLE shops
    ADD COLUMN IF NOT EXISTS last_heavy_sync_at       TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS cards_cursor_updated_at  TEXT,
    ADD COLUMN IF NOT EXISTS cards_cursor_nmid        BIGINT,
    ADD COLUMN IF NOT EXISTS last_full_cards_sync_at  TIMESTAMPTZ;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE shops
    DROP COLUMN IF EXISTS last_heavy_sync_at,
    DROP COLUMN IF EXISTS cards_cursor_updated_at,
    DROP COLUMN IF EXISTS cards_cursor_nmid,
    DROP COLUMN IF EXISTS last_full_cards_sync_at;

-- +goose StatementEnd
