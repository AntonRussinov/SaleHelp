-- +goose Up
-- +goose StatementBegin

CREATE TABLE dbs_orders (
    id               BIGSERIAL PRIMARY KEY,
    shop_id           BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    wb_order_id       BIGINT NOT NULL,
    order_uid         TEXT NOT NULL,
    nm_id             BIGINT NOT NULL,
    chrt_id           BIGINT NOT NULL,
    article           TEXT,
    supplier_status   TEXT,
    wb_status         TEXT,
    created_at        TIMESTAMPTZ NOT NULL,
    notified_at       TIMESTAMPTZ,
    extra             JSONB,

    UNIQUE (shop_id, wb_order_id)
);

CREATE INDEX idx_dbs_orders_shop_id ON dbs_orders(shop_id);
CREATE INDEX idx_dbs_orders_shop_notified ON dbs_orders(shop_id, notified_at)
    WHERE notified_at IS NULL;
CREATE INDEX idx_dbs_orders_created_at ON dbs_orders(shop_id, created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS dbs_orders CASCADE;

-- +goose StatementEnd
