-- +goose Up

CREATE TABLE ai_settings (
    id                  BIGSERIAL PRIMARY KEY,
    shop_id             BIGINT NOT NULL UNIQUE REFERENCES shops(id) ON DELETE CASCADE,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active           BOOLEAN NOT NULL DEFAULT false,
    tone_of_voice       TEXT NOT NULL DEFAULT 'friendly',
    response_length     TEXT NOT NULL DEFAULT 'medium', -- short | medium | long
    addressing          TEXT NOT NULL DEFAULT 'vy',     -- ty | vy
    emoji_mode          TEXT NOT NULL DEFAULT 'minimal', -- none | minimal | expressive
    brand_name          TEXT,
    brand_description   TEXT,
    representative_name TEXT,
    signature           TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_settings_user_id ON ai_settings(user_id);
CREATE INDEX idx_ai_settings_active ON ai_settings(shop_id) WHERE is_active = true;

CREATE TABLE ai_rules (
    id              BIGSERIAL PRIMARY KEY,
    ai_settings_id  BIGINT NOT NULL REFERENCES ai_settings(id) ON DELETE CASCADE,
    rule_type       TEXT NOT NULL CHECK (rule_type IN ('rule', 'forbidden_word')),
    rule_text       TEXT NOT NULL,
    order_index     INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_rules_settings_id ON ai_rules(ai_settings_id);

CREATE TABLE ai_templates (
    id                   BIGSERIAL PRIMARY KEY,
    ai_settings_id       BIGINT NOT NULL REFERENCES ai_settings(id) ON DELETE CASCADE,
    category             TEXT NOT NULL CHECK (category IN ('positive', 'neutral', 'negative')),
    custom_prompt        TEXT,
    tone_voice_override  TEXT,
    enabled              BOOLEAN NOT NULL DEFAULT true,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (ai_settings_id, category)
);

-- ai_examples: few-shot примеры (Pro, таблица создана, в промпт пока не включается)
CREATE TABLE ai_examples (
    id              BIGSERIAL PRIMARY KEY,
    ai_settings_id  BIGINT NOT NULL REFERENCES ai_settings(id) ON DELETE CASCADE,
    review_text     TEXT NOT NULL,
    ideal_response  TEXT NOT NULL,
    category        TEXT CHECK (category IN ('positive', 'neutral', 'negative')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_examples_settings_id ON ai_examples(ai_settings_id);

-- Миграция существующих данных из shop_ai_settings
INSERT INTO ai_settings (shop_id, user_id, is_active, tone_of_voice, brand_name)
SELECT s.shop_id, sh.user_id, s.is_active, s.tone_of_voice, COALESCE(s.brand_name, '')
FROM shop_ai_settings s
JOIN shops sh ON sh.id = s.shop_id
ON CONFLICT (shop_id) DO NOTHING;

DROP TABLE IF EXISTS shop_ai_settings;

-- +goose Down
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

DROP TABLE IF EXISTS ai_examples;
DROP TABLE IF EXISTS ai_templates;
DROP TABLE IF EXISTS ai_rules;
DROP TABLE IF EXISTS ai_settings;
