package domain

import (
	"encoding/json"
	"time"
)

// TelegramChatState — состояние диалога (wizard) для Telegram чата.
type TelegramChatState struct {
	ChatID    int64
	State     string          // idle | wait_wb_token | wait_shop_name | wait_invite_token
	Payload   json.RawMessage // JSONB
	UpdatedAt time.Time
}
