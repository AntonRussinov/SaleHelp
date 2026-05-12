package model

import (
	"salehelp/domain"
	"time"
)

type AIPackage struct {
	ID          int       `gorm:"primaryKey;column:id"`
	UserID      int       `gorm:"column:user_id;not null"`
	Size        int       `gorm:"column:size;not null"`
	Remaining   int       `gorm:"column:remaining;not null"`
	PriceMinor  int       `gorm:"column:price_minor;not null"`
	PurchasedAt time.Time `gorm:"column:purchased_at;not null"`
}

func (*AIPackage) TableName() string { return "ai_packages" }

func (p *AIPackage) ToDomain() *domain.AIPackage {
	if p == nil {
		return nil
	}
	return &domain.AIPackage{
		ID:          p.ID,
		UserID:      p.UserID,
		Size:        p.Size,
		Remaining:   p.Remaining,
		PriceMinor:  p.PriceMinor,
		PurchasedAt: p.PurchasedAt,
	}
}
