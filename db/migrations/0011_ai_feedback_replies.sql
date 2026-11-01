-- +goose Up
-- +goose StatementBegin

-- 1.1 Настройки AI для магазина
CREATE TABLE shop_ai_settings (
    shop_id       BIGINT PRIMARY KEY REFERENCES shops(id) ON DELETE CASCADE,
    tone_of_voice TEXT NOT NULL DEFAULT 'standard'
        CHECK (tone_of_voice IN ('standard', 'friendly', 'official', 'empathic', 'minimalist')),
    brand_name    TEXT,
    custom_rules  TEXT,
    is_active     BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_shop_ai_settings_is_active ON shop_ai_settings(is_active) WHERE is_active = true;

-- 1.2 Изменение wb_feedbacks: is_replied, replied_at
ALTER TABLE wb_feedbacks
    ADD COLUMN IF NOT EXISTS is_replied BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS replied_at TIMESTAMPTZ;

-- 1.3 Логи ответов AI (допускаем несколько попыток на отзыв)
CREATE TABLE ai_feedback_replies (
    id            BIGSERIAL PRIMARY KEY,
    feedback_id   BIGINT NOT NULL REFERENCES wb_feedbacks(id) ON DELETE CASCADE,
    shop_id       BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    answer_text   TEXT,
    model_used    TEXT NOT NULL,
    status        TEXT NOT NULL
        CHECK (status IN ('generated', 'error', 'skipped', 'no_quota')),
    error_message TEXT,
    tokens_used   INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_feedback_replies_feedback_created ON ai_feedback_replies(feedback_id, created_at DESC);
CREATE INDEX idx_ai_feedback_replies_shop_created ON ai_feedback_replies(shop_id, created_at DESC);

-- 2.1 Баланс/квоты ответов по магазину и месяцу
CREATE TABLE shop_ai_reply_quotas (
    id             BIGSERIAL PRIMARY KEY,
    shop_id        BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    month          DATE NOT NULL,
    included_limit INT NOT NULL,
    included_used  INT NOT NULL DEFAULT 0,
    paid_balance   INT NOT NULL DEFAULT 0,
    paid_used      INT NOT NULL DEFAULT 0,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (shop_id, month)
);

CREATE INDEX idx_shop_ai_reply_quotas_shop_month ON shop_ai_reply_quotas(shop_id, month);

-- 2.2 Покупки пакетов ответов
CREATE TABLE shop_ai_reply_purchases (
    id            BIGSERIAL PRIMARY KEY,
    shop_id       BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    package_name  TEXT NOT NULL,
    answers_count INT NOT NULL,
    price_rub     INT NOT NULL,
    payment_ref   TEXT,
    status        TEXT NOT NULL
        CHECK (status IN ('pending', 'paid', 'canceled')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    paid_at       TIMESTAMPTZ
);

CREATE INDEX idx_shop_ai_reply_purchases_shop ON shop_ai_reply_purchases(shop_id);

-- 2.3 Аудит списаний
CREATE TABLE shop_ai_reply_ledger (
    id          BIGSERIAL PRIMARY KEY,
    shop_id     BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    feedback_id BIGINT REFERENCES wb_feedbacks(id) ON DELETE SET NULL,
    month       DATE NOT NULL,
    source      TEXT NOT NULL CHECK (source IN ('included', 'paid')),
    delta       INT NOT NULL,
    reason      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_shop_ai_reply_ledger_shop_month ON shop_ai_reply_ledger(shop_id, month);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS shop_ai_reply_ledger;
DROP TABLE IF EXISTS shop_ai_reply_purchases;
DROP TABLE IF EXISTS shop_ai_reply_quotas;
DROP TABLE IF EXISTS ai_feedback_replies;

ALTER TABLE wb_feedbacks
    DROP COLUMN IF EXISTS is_replied,
    DROP COLUMN IF EXISTS replied_at;

DROP TABLE IF EXISTS shop_ai_settings;

-- +goose StatementEnd
