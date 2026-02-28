package model

import (
	"salehelp/domain"
	"time"
)

type User struct {
	ID             int       `gorm:"primaryKey;column:id"`
	Email          *string   `gorm:"column:email;type:text;unique"`
	PasswordHash   *string   `gorm:"column:password_hash;type:text"`
	GoogleID       *string   `gorm:"column:google_id;type:text;unique"`
	TelegramChatID *int64    `gorm:"column:telegram_chat_id;type:bigint;unique"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func UserColumns() []string {
	return []string{
		"id",
		"email",
		"password_hash",
		"google_id",
		"telegram_chat_id",
		"created_at",
		"updated_at",
	}
}

func (u *User) ToDomain() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{
		ID:             u.ID,
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		GoogleID:       u.GoogleID,
		TelegramChatID: u.TelegramChatID,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func (*User) TableName() string {
	return "users"
}
