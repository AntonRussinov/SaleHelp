-- +goose Up
-- +goose StatementBegin

ALTER TABLE wb_feedbacks
    ADD COLUMN IF NOT EXISTS order_status TEXT NOT NULL DEFAULT '';

COMMENT ON COLUMN wb_feedbacks.order_status IS 'Статус заказа из WB API: buyout — выкуп, returned — товар возвращён (для AI-ответа: не предлагать возврат, если returned)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE wb_feedbacks DROP COLUMN IF EXISTS order_status;

-- +goose StatementEnd
