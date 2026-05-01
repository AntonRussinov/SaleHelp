-- +goose Up
ALTER TABLE ai_feedback_replies ADD COLUMN IF NOT EXISTS wb_sent_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE ai_feedback_replies DROP COLUMN IF EXISTS wb_sent_at;
