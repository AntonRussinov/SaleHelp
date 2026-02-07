-- +goose Up
-- +goose StatementBegin

ALTER TABLE shops
  ADD COLUMN IF NOT EXISTS wb_report_timezone TEXT;

COMMENT ON COLUMN shops.wb_report_timezone IS 'Часовой пояс для отчётов WB (sales, orders), IANA: Europe/Moscow, Asia/Almaty и т.д. Если NULL — используется глобальный из конфига.';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE shops DROP COLUMN IF EXISTS wb_report_timezone;

-- +goose StatementEnd
