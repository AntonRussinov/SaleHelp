package domain

import "time"

// ShopInviteToken — инвайт для подключения нового Telegram-чата к существующему магазину.
type ShopInviteToken struct {
	ID              int
	ShopID          int
	Token           string
	Role            string // "member" | "owner" (на будущее)
	ExpiresAt       time.Time
	UsedAt          *time.Time
	CreatedByChatID int64
	CreatedAt       time.Time
}
