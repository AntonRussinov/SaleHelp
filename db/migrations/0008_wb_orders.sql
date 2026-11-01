-- +goose Up
-- +goose StatementBegin

-- Все заказы из WB Orders API (/api/v1/supplier/orders).
-- Нужны для раннего индикатора спроса в stock_replenish: orders (+1) и отмены/отказы is_cancel (+-1) по cancel_date.
-- Связь с sales по srid.
CREATE TABLE wb_orders (
    id                BIGSERIAL PRIMARY KEY,
    shop_id           BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    srid              TEXT NOT NULL,
    nm_id             BIGINT NOT NULL,
    barcode           TEXT,
    supplier_article  TEXT,
    date              TIMESTAMPTZ NOT NULL,
    last_change_date  TIMESTAMPTZ NOT NULL,
    is_cancel         BOOLEAN NOT NULL DEFAULT FALSE,
    cancel_date       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL,

    UNIQUE (shop_id, srid)
);

CREATE INDEX idx_wb_orders_shop_id ON wb_orders(shop_id);
CREATE INDEX idx_wb_orders_shop_date ON wb_orders(shop_id, date DESC);
CREATE INDEX idx_wb_orders_shop_cancel_date ON wb_orders(shop_id, cancel_date DESC) WHERE is_cancel = true;
CREATE INDEX idx_wb_orders_shop_last_change ON wb_orders(shop_id, last_change_date DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wb_orders;

-- +goose StatementEnd

