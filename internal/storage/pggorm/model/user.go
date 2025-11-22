package model

import (
	"salehelp/domain"
	"time"
)

type User struct {
	ID           int       `gorm:"primaryKey;column:id"`
	Email        string    `gorm:"column:email;type:text;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func UserColumns() []string {
	return []string{
		"id",
		"email",
		"password_hash",
		"created_at",
		"updated_at",
	}
}

func (u *User) ToDomain() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (*User) TableName() string {
	return "users"
}
