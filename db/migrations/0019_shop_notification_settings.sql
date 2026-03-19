-- +goose Up
-- +goose StatementBegin

-- Персональные настройки уведомлений по магазину (включено/выключено по типу события).
CREATE TABLE IF NOT EXISTS shop_notification_settings (
    id         BIGSERIAL PRIMARY KEY,
    shop_id    BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    event_key  TEXT NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (shop_id, event_key)
);

CREATE INDEX IF NOT EXISTS idx_shop_notification_settings_shop_id
    ON shop_notification_settings (shop_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS shop_notification_settings CASCADE;

-- +goose StatementEnd

