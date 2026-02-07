-- +goose Up
-- +goose StatementBegin

CREATE TABLE wb_feedbacks (
    id                BIGSERIAL PRIMARY KEY,
    shop_id           BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    wb_feedback_id    TEXT NOT NULL,
    nm_id             BIGINT NOT NULL,
    product_name      TEXT NOT NULL DEFAULT '',
    text              TEXT NOT NULL DEFAULT '',
    pros              TEXT NOT NULL DEFAULT '',
    cons              TEXT NOT NULL DEFAULT '',
    product_valuation INT NOT NULL DEFAULT 0,
    user_name         TEXT NOT NULL DEFAULT '',
    created_at        TIMESTAMPTZ NOT NULL,
    notified_at       TIMESTAMPTZ,
    extra             JSONB,

    UNIQUE (shop_id, wb_feedback_id)
);

CREATE INDEX idx_wb_feedbacks_shop_id ON wb_feedbacks(shop_id);
CREATE INDEX idx_wb_feedbacks_shop_notified ON wb_feedbacks(shop_id, notified_at)
    WHERE notified_at IS NULL;
CREATE INDEX idx_wb_feedbacks_created_at ON wb_feedbacks(shop_id, created_at DESC);

CREATE TABLE wb_questions (
    id             BIGSERIAL PRIMARY KEY,
    shop_id        BIGINT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    wb_question_id TEXT NOT NULL,
    nm_id          BIGINT NOT NULL,
    product_name   TEXT NOT NULL DEFAULT '',
    text           TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL,
    notified_at    TIMESTAMPTZ,
    extra          JSONB,

    UNIQUE (shop_id, wb_question_id)
);

CREATE INDEX idx_wb_questions_shop_id ON wb_questions(shop_id);
CREATE INDEX idx_wb_questions_shop_notified ON wb_questions(shop_id, notified_at)
    WHERE notified_at IS NULL;
CREATE INDEX idx_wb_questions_created_at ON wb_questions(shop_id, created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wb_questions CASCADE;
DROP TABLE IF EXISTS wb_feedbacks CASCADE;

-- +goose StatementEnd
