package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID               int            `gorm:"primaryKey;column:id"`
	MasterSKU        string         `gorm:"column:master_sku;type:varchar(50);unique;not null"`
	Model            string         `gorm:"column:model;type:varchar(255);not null"`
	Brand            string         `gorm:"column:brand;type:varchar(100)"`
	MasterTitle      string         `gorm:"column:master_title;type:varchar(255)"`
	Title            string         `gorm:"column:title;type:varchar(255)"`
	VerticalCategory string         `gorm:"column:vertical_category;type:varchar(100)"`
	MasterCategory   string         `gorm:"column:master_category;type:varchar(100)"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func ProductColumns() []string {
	return []string{
		"id",
		"master_sku",
		"model",
		"brand",
		"master_title",
		"title",
		"vertical_category",
		"master_category",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (p *Product) ToDomain() *domain.Product {
	if p == nil {
		return nil
	}
	return &domain.Product{
		ID:               p.ID,
		MasterSKU:        p.MasterSKU,
		Model:            p.Model,
		Brand:            p.Brand,
		MasterTitle:      p.MasterTitle,
		Title:            p.Title,
		VerticalCategory: p.VerticalCategory,
		MasterCategory:   p.MasterCategory,
	}
}

func (*Product) TableName() string {
	return "products"
}
