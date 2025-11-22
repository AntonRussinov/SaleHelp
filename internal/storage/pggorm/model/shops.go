package model

import (
	"salehelp/domain"
	"time"
)

type Shop struct {
	ID          int       `gorm:"primaryKey;column:id"`
	UserID      int       `gorm:"column:user_id;not null"`
	Marketplace string    `gorm:"column:marketplace;type:text;not null"` // "WB", "OZON", "Kaspi", ...
	ExternalID  string    `gorm:"column:external_id;type:text"`
	Name        string    `gorm:"column:name;type:text;not null"`
	APIToken    string    `gorm:"column:api_token;type:text;not null"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func ShopColumns() []string {
	return []string{
		"id",
		"user_id",
		"marketplace",
		"external_id",
		"name",
		"api_token",
		"is_active",
		"created_at",
		"updated_at",
	}
}

func (s *Shop) ToDomain() *domain.Shop {
	if s == nil {
		return nil
	}

	return &domain.Shop{
		ID:          s.ID,
		UserID:      s.UserID,
		Marketplace: s.Marketplace,
		ExternalID:  s.ExternalID,
		Name:        s.Name,
		APIToken:    s.APIToken,
		IsActive:    s.IsActive,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (*Shop) TableName() string {
	return "shops"
}
