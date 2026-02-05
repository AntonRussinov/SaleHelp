-- +goose Up
-- +goose StatementBegin

-- Отказы (заказы не забраны): из WB Orders API (isCancel = true). Связь с sale_events по srid.
CREATE TABLE order_refusals (
    id                BIGSERIAL PRIMARY KEY,
    shop_id           BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    srid              TEXT NOT NULL,
    nm_id             BIGINT NOT NULL,
    barcode           TEXT,
    supplier_article  TEXT,
    date              TIMESTAMPTZ NOT NULL,
    cancel_date       TIMESTAMPTZ NOT NULL,
    last_change_date  TIMESTAMPTZ NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL,
    notified_at       TIMESTAMPTZ,

    UNIQUE (shop_id, srid)
);

CREATE INDEX idx_order_refusals_shop_id ON order_refusals(shop_id);
CREATE INDEX idx_order_refusals_shop_notified ON order_refusals(shop_id, notified_at)
    WHERE notified_at IS NULL;
CREATE INDEX idx_order_refusals_shop_last_change ON order_refusals(shop_id, last_change_date DESC);

-- Курсор пагинации Orders API по магазину (lastChangeDate последней обработанной строки).
CREATE TABLE order_sync_cursor (
    shop_id           BIGINT NOT NULL PRIMARY KEY REFERENCES shops(id) ON DELETE CASCADE,
    last_change_date  TIMESTAMPTZ NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_sync_cursor;
DROP TABLE IF EXISTS order_refusals;

-- +goose StatementEnd
