package model

import (
	"time"

	"salehelp/domain"
)

type ProductSizeSKU struct {
	ID            int       `gorm:"primaryKey;column:id"`
	ProductSizeID int       `gorm:"column:product_size_id;not null"`
	SKU           string    `gorm:"column:sku;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
}

func ProductSizeSKUColumns() []string {
	return []string{
		"id",
		"product_size_id",
		"sku",
		"created_at",
	}
}

func (p *ProductSizeSKU) ToDomain() *domain.ProductSizeSKU {
	if p == nil {
		return nil
	}

	return &domain.ProductSizeSKU{
		ID:            p.ID,
		ProductSizeID: p.ProductSizeID,
		SKU:           p.SKU,
		CreatedAt:     p.CreatedAt,
	}
}

func (*ProductSizeSKU) TableName() string {
	return "product_size_skus"
}
