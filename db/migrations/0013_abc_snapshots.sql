-- +goose Up
-- +goose StatementBegin

CREATE TABLE abc_snapshots (
    id           BIGSERIAL PRIMARY KEY,
    shop_id      BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    product_id   BIGINT REFERENCES products(id) ON DELETE SET NULL,
    nm_id        BIGINT NOT NULL,
    category     CHAR(1) NOT NULL CHECK (category IN ('A', 'B', 'C')),
    revenue      NUMERIC(14,2) NOT NULL DEFAULT 0,
    sales_count  INTEGER NOT NULL DEFAULT 0,
    period_days  SMALLINT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end   TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_abc_snapshots_shop_period
    ON abc_snapshots (shop_id, period_days, period_end DESC);

CREATE INDEX idx_abc_snapshots_shop_period_created
    ON abc_snapshots (shop_id, period_days, created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS abc_snapshots CASCADE;

-- +goose StatementEnd
