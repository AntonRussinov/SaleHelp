package model

import (
	"time"

	"salehelp/domain"
)

type ProductSize struct {
	ID           int       `gorm:"primaryKey;column:id"`
	ProductID    int       `gorm:"column:product_id;not null"`
	ExternalID   string    `gorm:"column:external_id;not null"`
	TechSizeName *string   `gorm:"column:tech_size_name"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

func ProductSizeColumns() []string {
	return []string{
		"id",
		"product_id",
		"external_id",
		"tech_size_name",
		"created_at",
		"updated_at",
	}
}

func (p *ProductSize) ToDomain() *domain.ProductSize {
	if p == nil {
		return nil
	}

	return &domain.ProductSize{
		ID:           p.ID,
		ProductID:    p.ProductID,
		ExternalID:   p.ExternalID,
		TechSizeName: p.TechSizeName,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func (*ProductSize) TableName() string {
	return "product_sizes"
}
