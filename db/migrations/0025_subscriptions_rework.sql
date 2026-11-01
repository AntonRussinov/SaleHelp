-- +goose Up

-- Drop old subscription system (data loss acceptable, only dev data)
DROP TABLE IF EXISTS subscription_notification_events;
DROP TABLE IF EXISTS subscription_addons;
DROP TABLE IF EXISTS subscriptions;

-- Subscriptions: user-based, with tier, AI quota, shop limits
CREATE TABLE subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier            TEXT NOT NULL,          -- start | business | pro | trial | package_only
    status          TEXT NOT NULL,          -- active | trial | cancelled | expired
    period_months   INT NOT NULL DEFAULT 1, -- 1, 3, 6, 12 (0 for trial)
    started_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL,
    trial_ends_at   TIMESTAMPTZ,            -- only for trial
    ai_quota        INT NOT NULL DEFAULT 0, -- AI replies limit this month
    ai_used         INT NOT NULL DEFAULT 0, -- AI replies used this month
    ai_resets_at    TIMESTAMPTZ NOT NULL,   -- start of next quota period
    max_shops       INT NOT NULL DEFAULT 1, -- shops included in plan
    extra_shops     INT NOT NULL DEFAULT 0, -- additional shops purchased
    price_minor     INT NOT NULL DEFAULT 0, -- monthly price in kopecks (0 for trial)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Only one active or trial subscription per user at a time
CREATE UNIQUE INDEX idx_subscriptions_active_user ON subscriptions(user_id)
    WHERE status IN ('active', 'trial');

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_expires_at ON subscriptions(expires_at);

-- AI packages: purchased bundles of additional AI replies
CREATE TABLE ai_packages (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    size         INT NOT NULL,          -- total replies in this package
    remaining    INT NOT NULL,          -- replies remaining
    price_minor  INT NOT NULL,          -- price paid in kopecks
    purchased_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_packages_user_remaining ON ai_packages(user_id, remaining)
    WHERE remaining > 0;

-- +goose Down
DROP TABLE IF EXISTS ai_packages;
DROP TABLE IF EXISTS subscriptions;
