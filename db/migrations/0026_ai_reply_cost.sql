-- +goose Up
ALTER TABLE ai_feedback_replies
    ADD COLUMN IF NOT EXISTS input_tokens          INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS output_tokens         INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ai_tier               TEXT NOT NULL DEFAULT 'paid',
    ADD COLUMN IF NOT EXISTS estimated_cost_kopeks INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE ai_feedback_replies
    DROP COLUMN IF EXISTS input_tokens,
    DROP COLUMN IF EXISTS output_tokens,
    DROP COLUMN IF EXISTS ai_tier,
    DROP COLUMN IF EXISTS estimated_cost_kopeks;
