-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

COMMENT ON COLUMN users.role IS 'Роль пользователя: user, admin и т.д.';

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS role;
