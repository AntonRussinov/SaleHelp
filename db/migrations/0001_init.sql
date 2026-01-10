-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id                  BIGSERIAL PRIMARY KEY,
    telegram_chat_id    BIGINT UNIQUE,
    email               TEXT UNIQUE,
    password_hash       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
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
    title           TEXT NOT NULL,
    url             TEXT,
    image_url       TEXT,
    brand           TEXT,
    category        TEXT,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (shop_id, external_id)
);

CREATE INDEX idx_products_shop_id ON products(shop_id);

CREATE TABLE product_snapshots (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    -- момент, на который снят слепок (может отличаться от created_at)
    snapshot_time   TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Цены и валюта по товару в целом (агрегаты по размерам)
    currency_code           CHAR(3),         -- 'RUB'
    base_price              NUMERIC(12,2),   -- WB: price (минимальный/средний по размерам)
    final_price             NUMERIC(12,2),   -- WB: discountedPrice
    club_price              NUMERIC(12,2),   -- WB: clubDiscountedPrice
    discount                NUMERIC(5,2),    -- % скидки
    club_discount           NUMERIC(5,2),    -- % клубной скидки
    old_price               NUMERIC(12,2),   -- при желании можно использовать для твоей логики
                                             -- (необязателен, но раз хочешь "максимум" — оставим)
    -- Остатки и продажи по всему товару
    stock_total             INTEGER,         -- суммарный остаток по всем размерам и складам
    stock_in_way            INTEGER,         -- если из stock API можно вытащить "в пути"

    rating                  NUMERIC(3,2),
    reviews_count           INTEGER,

    sales_24h               INTEGER,
    sales_7d                INTEGER,

    position_search         INTEGER,         -- позиция в поиске по основному ключевому

    -- Флаги с уровня товара (из API типа isBadTurnover и пр.)
    is_bad_turnover         BOOLEAN,
    editable_size_price     BOOLEAN,

    -- Запасной карман под всё, что не хочется сейчас раскладывать по колонкам
    extra                   JSONB,

    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- полезные индексы
CREATE INDEX idx_product_snapshots_product_time
    ON product_snapshots (product_id, snapshot_time DESC);

CREATE INDEX IF NOT EXISTS idx_product_snapshots_product_time_id
ON product_snapshots (product_id, snapshot_time DESC, id DESC);


CREATE TABLE product_sizes (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    external_id     TEXT NOT NULL,   -- chrtID (TEXT для мульти-маркетплейса)
    tech_size_name  TEXT,            -- "M", "42" и т.д.

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (product_id, external_id)
);

CREATE INDEX idx_product_sizes_product_id
    ON product_sizes(product_id);

CREATE TABLE product_size_snapshots (
    id                  BIGSERIAL PRIMARY KEY,

    product_snapshot_id BIGINT NOT NULL
        REFERENCES product_snapshots(id) ON DELETE CASCADE,

    product_size_id     BIGINT NOT NULL
        REFERENCES product_sizes(id) ON DELETE CASCADE,

    sku                 TEXT NOT NULL,   -- SKU, по которому пришёл stock
    stock               INTEGER,

    stock_by_warehouse  JSONB,
    extra               JSONB,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (product_snapshot_id, product_size_id, sku)
);

CREATE INDEX idx_pss_product_size
    ON product_size_snapshots(product_size_id);

CREATE INDEX idx_pss_sku
    ON product_size_snapshots(sku);

CREATE INDEX IF NOT EXISTS idx_pss_size_created_id
    ON product_size_snapshots (product_size_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_pss_snapshot_size
    ON product_size_snapshots (product_snapshot_id, product_size_id);



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
    scope               TEXT NOT NULL DEFAULT 'auto' CHECK (scope IN ('auto','product','size')),
    product_size_id     BIGINT REFERENCES product_sizes(id) ON DELETE CASCADE,
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
  id                BIGSERIAL PRIMARY KEY,
  shop_id           BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,

  plan              TEXT NOT NULL DEFAULT 'basic',
  status            TEXT NOT NULL, -- active | expired | canceled

  valid_from        TIMESTAMPTZ NOT NULL,
  valid_to          TIMESTAMPTZ NOT NULL,

  period_months     SMALLINT NOT NULL,

  base_amount_minor BIGINT NOT NULL,
  currency_code     CHAR(3) NOT NULL,

  paid_at           TIMESTAMPTZ,
  source            TEXT NOT NULL DEFAULT 'manual',
  comment           TEXT,

  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

  CHECK (valid_to > valid_from),
  CHECK (period_months IN (1,3,6)),
  CHECK (status IN ('active','expired','canceled'))
);

CREATE UNIQUE INDEX ux_subscriptions_one_active_per_shop
  ON subscriptions (shop_id)
  WHERE status = 'active';

CREATE INDEX idx_subscriptions_shop_status_validto
  ON subscriptions (shop_id, status, valid_to DESC);


CREATE TABLE subscription_addons (
  id               BIGSERIAL PRIMARY KEY,
  subscription_id  BIGINT NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,

  code             TEXT NOT NULL, -- telegram_extra_chat, ...
  quantity         INT  NOT NULL DEFAULT 1,

  amount_minor     BIGINT NOT NULL DEFAULT 0,
  currency_code    CHAR(3) NOT NULL,

  meta             JSONB NOT NULL DEFAULT '{}'::jsonb,

  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

  CHECK (quantity >= 0),
  CHECK (code <> '')
);

CREATE UNIQUE INDEX ux_addons_subscription_code
  ON subscription_addons (subscription_id, code);

CREATE INDEX idx_addons_subscription
  ON subscription_addons (subscription_id);

CREATE TABLE alerts_triggered (
    id              BIGSERIAL PRIMARY KEY,
    alert_id        BIGINT NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    snapshot_id     BIGINT NOT NULL REFERENCES product_snapshots(id) ON DELETE CASCADE,
    shop_id         BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    message         TEXT,
    product_size_id BIGINT REFERENCES product_sizes(id) ON DELETE CASCADE,
    alert_type      TEXT   NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_alerts_triggered_alert ON alerts_triggered(alert_id);
CREATE INDEX idx_alerts_triggered_product ON alerts_triggered(product_id);
CREATE INDEX idx_alerts_triggered_snapshot ON alerts_triggered(snapshot_id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_alerts_triggered_once
    ON alerts_triggered(alert_id, product_id, product_size_id, snapshot_id);
CREATE INDEX IF NOT EXISTS idx_alerts_triggered_shop_created
    ON alerts_triggered (shop_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_alerts_triggered_product_created
    ON alerts_triggered (product_id, created_at DESC);


CREATE TABLE shop_notification_endpoints (
  id           BIGSERIAL PRIMARY KEY,
  shop_id      BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  channel      TEXT   NOT NULL,                -- telegram | email | whatsapp
  address      TEXT   NOT NULL,                -- chat_id / email / phone
  enabled      BOOLEAN NOT NULL DEFAULT TRUE,
  label        TEXT,                           -- "CEO", "Warehouse", "Sales", ...
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (shop_id, channel, address)
);

CREATE INDEX idx_endpoints_shop_enabled ON shop_notification_endpoints(shop_id, enabled);
CREATE INDEX idx_endpoints_channel ON shop_notification_endpoints(channel);
CREATE INDEX IF NOT EXISTS idx_endpoints_shop_channel_enabled
ON shop_notification_endpoints (shop_id, channel, enabled);


CREATE TABLE alerts_triggered_deliveries (
  id           BIGSERIAL PRIMARY KEY,
  triggered_id BIGINT NOT NULL REFERENCES alerts_triggered(id) ON DELETE CASCADE,
  endpoint_id  BIGINT NOT NULL REFERENCES shop_notification_endpoints(id) ON DELETE CASCADE,

  status       TEXT NOT NULL DEFAULT 'new',     -- new | sent | failed
  sent_at      TIMESTAMPTZ,
  attempts     INT NOT NULL DEFAULT 0,
  last_error   TEXT,

  lease_until  TIMESTAMPTZ,
  leased_by    TEXT,

  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (triggered_id, endpoint_id)
);

CREATE INDEX idx_deliveries_ready
ON alerts_triggered_deliveries (lease_until, created_at)
WHERE status='new';
CREATE INDEX IF NOT EXISTS idx_deliveries_pick
ON alerts_triggered_deliveries (status, lease_until, created_at, id);


CREATE TABLE wb_warehouses (
    id               BIGSERIAL PRIMARY KEY,
    wb_id            BIGINT NOT NULL,         -- id из WB
    name             TEXT NOT NULL,
    address          TEXT,
    city             TEXT,
    longitude        NUMERIC(10,6),
    latitude         NUMERIC(10,6),
    cargo_type       INTEGER,
    delivery_type    INTEGER,
    federal_district TEXT,
    selected         BOOLEAN,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (wb_id)
);

CREATE INDEX idx_wb_warehouses_wb_id ON wb_warehouses(wb_id);

-- Склады продавца (привязаны к shop)
CREATE TABLE seller_warehouses (
    id             BIGSERIAL PRIMARY KEY,
    shop_id        BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    wb_office_id   BIGINT NOT NULL,           -- officeId из WB
    wb_id          BIGINT NOT NULL,           -- id склада из WB
    name           TEXT NOT NULL,
    cargo_type     INTEGER,
    delivery_type  INTEGER,
    is_deleting    BOOLEAN NOT NULL DEFAULT FALSE,
    is_processing  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (shop_id, wb_office_id)
);

CREATE INDEX idx_seller_warehouses_shop ON seller_warehouses(shop_id);


CREATE TABLE product_size_skus (
    id              BIGSERIAL PRIMARY KEY,
    product_size_id BIGINT NOT NULL
        REFERENCES product_sizes(id) ON DELETE CASCADE,

    sku             TEXT NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (product_size_id, sku)
);

CREATE INDEX idx_product_size_skus_sku
    ON product_size_skus(sku);


CREATE TABLE product_size_price_snapshots (
    id                  BIGSERIAL PRIMARY KEY,
    product_snapshot_id BIGINT NOT NULL REFERENCES product_snapshots(id) ON DELETE CASCADE,
    product_size_id     BIGINT NOT NULL REFERENCES product_sizes(id) ON DELETE CASCADE,

    base_price          NUMERIC(12,2),
    final_price         NUMERIC(12,2),
    club_price          NUMERIC(12,2),
    discount            NUMERIC(5,2),
    club_discount       NUMERIC(5,2),

    extra               JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (product_snapshot_id, product_size_id)
);

CREATE INDEX idx_psps_snapshot ON product_size_price_snapshots(product_snapshot_id);
CREATE INDEX idx_psps_size ON product_size_price_snapshots(product_size_id);
CREATE INDEX IF NOT EXISTS idx_psps_size_snapshot_id
ON product_size_price_snapshots (product_size_id, product_snapshot_id DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_psps_snapshot_size
ON product_size_price_snapshots (product_snapshot_id, product_size_id);


CREATE TABLE shop_sync_leases (
  shop_id      BIGINT PRIMARY KEY REFERENCES shops(id) ON DELETE CASCADE,
  leased_by    TEXT,
  lease_until  TIMESTAMPTZ
);

CREATE INDEX idx_shop_sync_leases_until ON shop_sync_leases(lease_until);


CREATE TABLE telegram_chat_states (
  chat_id     BIGINT PRIMARY KEY,
  state       TEXT NOT NULL,      -- 'idle' | 'wait_wb_token' | 'wait_shop_name'
  payload     JSONB NOT NULL DEFAULT '{}'::jsonb,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_telegram_chat_states_updated
  ON telegram_chat_states(updated_at DESC);


CREATE TABLE shop_invite_tokens (
  id          BIGSERIAL PRIMARY KEY,
  shop_id     BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  token       TEXT NOT NULL UNIQUE,
  role        TEXT NOT NULL DEFAULT 'member', -- member | owner (на будущее)
  expires_at  TIMESTAMPTZ NOT NULL,
  used_at     TIMESTAMPTZ,
  created_by_chat_id BIGINT NOT NULL, -- кто создал (owner chat)
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_shop_invite_tokens_shop
  ON shop_invite_tokens (shop_id, created_at DESC);

CREATE INDEX idx_shop_invite_tokens_valid
  ON shop_invite_tokens (token)
  WHERE used_at IS NULL;


-- Вставляем временно тестового пользователя и магазин для отладки
INSERT INTO users (
    id,
    email,
    password_hash,
    created_at,
    updated_at
) VALUES (
    1,
    'test@example.com',
    'hash', 
    now(),
    now()
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS alerts_triggered_deliveries;
DROP TABLE IF EXISTS shop_notification_endpoints;

-- 1. Сначала зависимая таблица alerts_triggered
DROP INDEX IF EXISTS idx_alerts_triggered_snapshot;
DROP INDEX IF EXISTS idx_alerts_triggered_product;
DROP INDEX IF EXISTS idx_alerts_triggered_alert;
DROP TABLE IF EXISTS alerts_triggered;

-- 2. Снимаем FK с alerts на ai_recommendations
ALTER TABLE IF EXISTS alerts DROP CONSTRAINT IF EXISTS fk_alerts_ai_rec;

-- 3. Дропаем ai_recommendations (она ссылается на alerts)
DROP TABLE IF EXISTS ai_recommendations;

-- 4. Дропаем alerts и связанные с ними индексы
DROP INDEX IF EXISTS idx_alerts_product;
DROP INDEX IF EXISTS idx_alerts_shop_status;
DROP TABLE IF EXISTS alerts;

-- 5. Сначала дочерняя таблица размеров, потом снапшоты товаров
DROP INDEX IF EXISTS idx_product_size_snapshots_snapshot;
DROP INDEX IF EXISTS idx_product_size_snapshots_size;
DROP TABLE IF EXISTS product_size_snapshots;
DROP TABLE IF EXISTS product_size_price_snapshots;
DROP INDEX IF EXISTS idx_product_snapshots_product_time;
DROP TABLE IF EXISTS product_snapshots;

DROP TABLE IF EXISTS product_size_skus;
DROP TABLE IF EXISTS product_sizes;



-- 6. Продукты и индекс
DROP INDEX IF EXISTS idx_products_shop_id;
DROP TABLE IF EXISTS products;

-- 7. Подписки и индекс
DROP INDEX IF EXISTS idx_subscriptions_user;
DROP TABLE IF EXISTS subscriptions;

-- 13. Дропаем склады WB и индекс
DROP TABLE IF EXISTS wb_warehouses;

-- 12. Дропаем склады и индекс
DROP TABLE IF EXISTS seller_warehouses;

-- 11. Дропаем индексы и таблицы размеров товаров, снапшотов товаров и самих товаров, если они остались
DROP INDEX IF EXISTS idx_product_size_snapshots_snapshot;
DROP INDEX IF EXISTS idx_product_size_snapshots_size;
DROP TABLE IF EXISTS product_size_snapshots;

-- 9. Магазины и индексы
DROP INDEX IF EXISTS idx_shops_marketplace_external;
DROP INDEX IF EXISTS idx_shops_user_id;
DROP TABLE IF EXISTS shops;

-- 10. Пользователи (последними, т.к. на них всё ссылается)
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS shop_sync_leases;
DROP TABLE IF EXISTS telegram_chat_states;
DROP TABLE IF EXISTS shop_invite_tokens;
DROP TABLE IF EXISTS subscription_addons;
-- +goose StatementEnd

 