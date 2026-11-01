package domain

import "time"

// LinkToken — одноразовый токен для привязки Telegram ↔ аккаунт на сайте.
// Либо создаётся на сайте (UserID задан, TelegramChatID nil) — код вводят в боте.
// Либо создаётся в боте (TelegramChatID задан, UserID nil) — код вводят на сайте.
type LinkToken struct {
	Token          string
	UserID         *int   // задан для кодов «с сайта»
	TelegramChatID *int64 // задан для кодов «из Telegram»
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
