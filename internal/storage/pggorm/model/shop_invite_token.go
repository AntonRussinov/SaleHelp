package model

import (
	"time"

	"salehelp/domain"
)

type ShopInviteToken struct {
	ID              int        `gorm:"column:id;primaryKey"`
	ShopID          int        `gorm:"column:shop_id;not null"`
	Token           string     `gorm:"column:token;not null;unique"`
	Role            string     `gorm:"column:role;not null"`
	ExpiresAt       time.Time  `gorm:"column:expires_at;not null"`
	UsedAt          *time.Time `gorm:"column:used_at"`
	CreatedByChatID int64      `gorm:"column:created_by_chat_id;not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null"`
}

func (ShopInviteToken) TableName() string { return "shop_invite_tokens" }

func (m ShopInviteToken) ToDomain() *domain.ShopInviteToken {
	return &domain.ShopInviteToken{
		ID:              m.ID,
		ShopID:          m.ShopID,
		Token:           m.Token,
		Role:            m.Role,
		ExpiresAt:       m.ExpiresAt,
		UsedAt:          m.UsedAt,
		CreatedByChatID: m.CreatedByChatID,
		CreatedAt:       m.CreatedAt,
	}
}
