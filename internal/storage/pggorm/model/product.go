package model

import (
	"salehelp/domain"
	"time"
)

type Product struct {
	ID          int       `gorm:"primaryKey;column:id"`
	ShopID      int       `gorm:"column:shop_id;not null"`
	ExternalID  string    `gorm:"column:external_id;type:text;not null"`
	Title       string    `gorm:"column:title;type:text;not null"`
	URL         string    `gorm:"column:url;type:text"`
	ImageURL    string    `gorm:"column:image_url;type:text"`
	Brand       string    `gorm:"column:brand;type:text"`
	Category    string    `gorm:"column:category;type:text"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func ProductColumns() []string {
	return []string{
		"id",
		"shop_id",
		"external_id",
		"title",
		"url",
		"image_url",
		"brand",
		"category",
		"description",
		"created_at",
		"updated_at",
	}
}

func (p *Product) ToDomain() *domain.Product {
	if p == nil {
		return nil
	}

	return &domain.Product{
		ID:          p.ID,
		ShopID:      p.ShopID,
		ExternalID:  p.ExternalID,
		Title:       p.Title,
		URL:         p.URL,
		ImageURL:    p.ImageURL,
		Brand:       p.Brand,
		Category:    p.Category,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (*Product) TableName() string {
	return "products"
}
