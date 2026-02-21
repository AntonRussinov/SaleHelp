-- +goose Up
-- +goose StatementBegin

CREATE TABLE shop_abc_settings (
    shop_id             BIGINT PRIMARY KEY REFERENCES shops(id) ON DELETE CASCADE,
    abc_weekly_enabled   BOOLEAN NOT NULL DEFAULT true,
    abc_monthly_enabled  BOOLEAN NOT NULL DEFAULT true,
    abc_alerts_enabled   BOOLEAN NOT NULL DEFAULT true,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS shop_abc_settings CASCADE;

-- +goose StatementEnd
