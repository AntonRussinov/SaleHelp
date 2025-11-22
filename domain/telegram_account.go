package domain

import "time"

type TelegramAccount struct {
	ID         int
	UserID     int
	TelegramID int64
	ChatID     int64
	Username   string
	CreatedAt  time.Time
}
