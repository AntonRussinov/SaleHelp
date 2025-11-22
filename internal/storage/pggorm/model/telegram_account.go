package model

import (
	"salehelp/domain"
	"time"
)

type TelegramAccount struct {
	ID         int       `gorm:"primaryKey;column:id"`
	UserID     int       `gorm:"column:user_id;not null"`
	TelegramID int64     `gorm:"column:telegram_id;not null"`
	ChatID     int64     `gorm:"column:chat_id;not null"`
	Username   string    `gorm:"column:username;type:text"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func TelegramAccountColumns() []string {
	return []string{
		"id",
		"user_id",
		"telegram_id",
		"chat_id",
		"username",
		"created_at",
	}
}

func (t *TelegramAccount) ToDomain() *domain.TelegramAccount {
	if t == nil {
		return nil
	}

	return &domain.TelegramAccount{
		ID:         t.ID,
		UserID:     t.UserID,
		TelegramID: t.TelegramID,
		ChatID:     t.ChatID,
		Username:   t.Username,
		CreatedAt:  t.CreatedAt,
	}
}

func (*TelegramAccount) TableName() string {
	return "telegram_accounts"
}
