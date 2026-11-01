-- +goose Up
-- +goose StatementBegin

CREATE TABLE link_tokens (
    token            TEXT PRIMARY KEY,
    user_id          BIGINT REFERENCES users(id) ON DELETE CASCADE,
    telegram_chat_id BIGINT,
    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_link_tokens_one_target CHECK (
        (user_id IS NOT NULL AND telegram_chat_id IS NULL) OR
        (user_id IS NULL AND telegram_chat_id IS NOT NULL)
    )
);

CREATE INDEX idx_link_tokens_expires_at ON link_tokens(expires_at);
CREATE INDEX idx_link_tokens_telegram_chat_id ON link_tokens(telegram_chat_id) WHERE telegram_chat_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS link_tokens;

-- +goose StatementEnd
