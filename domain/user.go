package domain

import (
	"time"
)

type (
	User struct {
		ID             int
		Email          *string
		PasswordHash   *string
		TelegramChatID *int64
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)
