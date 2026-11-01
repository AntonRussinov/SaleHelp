-- +goose Up
-- +goose StatementBegin

-- Пакет восстановления цен по нажатию «Вернуть цены назад» в алерте price_drop.
CREATE TABLE price_restore_batches (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    restored_at TIMESTAMPTZ
);

CREATE INDEX idx_price_restore_batches_restored
    ON price_restore_batches (restored_at)
    WHERE restored_at IS NULL;

-- Связь пакета с триггерами алертов (только price_drop), входящими в одно сообщение.
CREATE TABLE price_restore_batch_triggered (
    batch_id     BIGINT NOT NULL REFERENCES price_restore_batches(id) ON DELETE CASCADE,
    triggered_id BIGINT NOT NULL REFERENCES alerts_triggered(id) ON DELETE CASCADE,
    PRIMARY KEY (batch_id, triggered_id)
);

CREATE INDEX idx_price_restore_batch_triggered_batch
    ON price_restore_batch_triggered (batch_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS price_restore_batch_triggered CASCADE;
DROP TABLE IF EXISTS price_restore_batches CASCADE;

-- +goose StatementEnd
