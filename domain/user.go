package domain

import (
	"time"
)

type (
	User struct {
		ID             int
		Email          *string
		PasswordHash   *string
		GoogleID       *string // sub из Google id_token (для входа через Google)
		TelegramChatID *int64
		Role           string // user, admin и т.д.
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)
