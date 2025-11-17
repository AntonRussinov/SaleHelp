-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE shops (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    marketplace     TEXT   NOT NULL,
    external_id     TEXT,
    name            TEXT NOT NULL,
    api_token       TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_shops_user_id ON shops(user_id);
CREATE INDEX idx_shops_marketplace_external ON shops(marketplace, external_id);

CREATE TABLE products (
    id              BIGSERIAL PRIMARY KEY,
    shop_id         BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    external_id     TEXT NOT NULL,
    sku             TEXT,
    title           TEXT NOT NULL,
    url             TEXT,
    image_url       TEXT,
    brand           TEXT,
    category        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (shop_id, external_id)
);

CREATE INDEX idx_products_shop_id ON products(shop_id);

CREATE TABLE product_snapshots (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    snapshot_time   TIMESTAMPTZ NOT NULL DEFAULT now(),

    price           NUMERIC(12,2),
    old_price       NUMERIC(12,2),
    discount        NUMERIC(5,2),
    stock_total     INTEGER,
    rating          NUMERIC(3,2),
    reviews_count   INTEGER,
    sales_24h       INTEGER,
    sales_7d        INTEGER,
    position_search INTEGER,
    extra           JSONB,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_snapshots_product_time
    ON product_snapshots(product_id, snapshot_time DESC);

CREATE TABLE alerts (
    id                  BIGSERIAL PRIMARY KEY,
    shop_id             BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    product_id          BIGINT REFERENCES products(id) ON DELETE CASCADE,
    type                TEXT   NOT NULL,
    severity            SMALLINT NOT NULL DEFAULT 1,
    status              TEXT NOT NULL DEFAULT 'new',
    message             TEXT,
    payload             JSONB,
    triggered_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at         TIMESTAMPTZ,
    ai_recommendation_id BIGINT
);

CREATE INDEX idx_alerts_shop_status ON alerts(shop_id, status, triggered_at DESC);
CREATE INDEX idx_alerts_product ON alerts(product_id);

CREATE TABLE ai_recommendations (
    id              BIGSERIAL PRIMARY KEY,
    shop_id         BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    product_id      BIGINT REFERENCES products(id) ON DELETE SET NULL,
    alert_id        BIGINT REFERENCES alerts(id) ON DELETE SET NULL,
    type            TEXT NOT NULL,
    model           TEXT,
    prompt          TEXT,
    content         TEXT NOT NULL,
    meta            JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE alerts
    ADD CONSTRAINT fk_alerts_ai_rec
    FOREIGN KEY (ai_recommendation_id)
    REFERENCES ai_recommendations(id)
    ON DELETE SET NULL;

CREATE TABLE subscriptions (
    id                      BIGSERIAL PRIMARY KEY,
    user_id                 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan                    TEXT   NOT NULL,
    status                  TEXT   NOT NULL,
    valid_from              TIMESTAMPTZ NOT NULL,
    valid_to                TIMESTAMPTZ NOT NULL,
    external_customer_id    TEXT,
    external_subscription_id TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id, status);

CREATE TABLE telegram_accounts (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    telegram_id     BIGINT NOT NULL,
    chat_id         BIGINT NOT NULL,
    username        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (telegram_id),
    UNIQUE (user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE telegram_accounts;
DROP INDEX idx_subscriptions_user;
DROP TABLE subscriptions;
ALTER TABLE alerts DROP CONSTRAINT fk_alerts_ai_rec;
DROP TABLE ai_recommendations;
DROP INDEX idx_alerts_product;
DROP INDEX idx_alerts_shop_status;
DROP TABLE alerts;
DROP INDEX idx_snapshots_product_time;
DROP TABLE product_snapshots;
DROP INDEX idx_products_shop_id;
DROP TABLE products;
DROP INDEX idx_shops_marketplace_external;
DROP INDEX idx_shops_user_id;
DROP TABLE shops;
DROP TABLE users;

-- +goose StatementEnd
