-- +goose Up
-- +goose StatementBegin

CREATE TYPE sale_event_type AS ENUM ('buyout', 'return', 'partial_return');

CREATE TABLE sale_events (
    id                BIGSERIAL PRIMARY KEY,
    shop_id           BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    event_type        sale_event_type NOT NULL,
    wb_sale_id        TEXT NOT NULL,
    srid              TEXT NOT NULL,
    date              TIMESTAMPTZ NOT NULL,
    last_change_date  TIMESTAMPTZ NOT NULL,
    nm_id             BIGINT NOT NULL,
    barcode           TEXT,
    supplier_article  TEXT,
    warehouse_name    TEXT,
    category          TEXT,
    subject           TEXT,
    brand             TEXT,
    tech_size         TEXT,
    income_id         BIGINT,
    total_price       NUMERIC(12,2),
    discount_percent  NUMERIC(5,2),
    finished_price    NUMERIC(12,2),
    price_with_disc   NUMERIC(12,2),
    for_pay           NUMERIC(12,2),
    created_at        TIMESTAMPTZ NOT NULL,
    notified_at       TIMESTAMPTZ,
    extra             JSONB,

    UNIQUE (shop_id, srid)
);

CREATE INDEX idx_sale_events_shop_id ON sale_events(shop_id);
CREATE INDEX idx_sale_events_shop_notified ON sale_events(shop_id, notified_at)
    WHERE notified_at IS NULL;
CREATE INDEX idx_sale_events_shop_created ON sale_events(shop_id, created_at DESC);
CREATE INDEX idx_sale_events_shop_last_change ON sale_events(shop_id, last_change_date DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS sale_events CASCADE;
DROP TYPE IF EXISTS sale_event_type;

-- +goose StatementEnd
