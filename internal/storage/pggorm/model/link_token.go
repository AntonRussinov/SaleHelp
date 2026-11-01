package model

import (
	"time"

	"salehelp/domain"
)

type LinkToken struct {
	Token          string     `gorm:"column:token;primaryKey;type:text"`
	UserID         *int       `gorm:"column:user_id"`
	TelegramChatID *int64     `gorm:"column:telegram_chat_id"`
	ExpiresAt      time.Time  `gorm:"column:expires_at;not null"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null"`
}

func (LinkToken) TableName() string { return "link_tokens" }

func (m LinkToken) ToDomain() *domain.LinkToken {
	return &domain.LinkToken{
		Token:          m.Token,
		UserID:         m.UserID,
		TelegramChatID: m.TelegramChatID,
		ExpiresAt:      m.ExpiresAt,
		CreatedAt:      m.CreatedAt,
	}
}
