package model

import (
	"salehelp/domain"
	"time"
)

type SellerWarehouse struct {
	ID           int       `gorm:"primaryKey;column:id"`
	ShopID       int       `gorm:"column:shop_id;not null"`
	WBOfficeID   int64     `gorm:"column:wb_office_id;not null"`
	WBID         int64     `gorm:"column:wb_id;not null"`
	Name         string    `gorm:"column:name;not null"`
	CargoType    *int      `gorm:"column:cargo_type"`
	DeliveryType *int      `gorm:"column:delivery_type"`
	IsDeleting   bool      `gorm:"column:is_deleting;not null;default:false"`
	IsProcessing bool      `gorm:"column:is_processing;not null;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func SellerWarehouseColumns() []string {
	return []string{
		"id",
		"shop_id",
		"wb_office_id",
		"wb_id",
		"name",
		"cargo_type",
		"delivery_type",
		"is_deleting",
		"is_processing",
		"created_at",
		"updated_at",
	}
}

func (m *SellerWarehouse) ToDomain() *domain.SellerWarehouse {
	return &domain.SellerWarehouse{
		ID:           m.ID,
		ShopID:       m.ShopID,
		WBOfficeID:   m.WBOfficeID,
		WBID:         m.WBID,
		Name:         m.Name,
		CargoType:    m.CargoType,
		DeliveryType: m.DeliveryType,
		IsDeleting:   m.IsDeleting,
		IsProcessing: m.IsProcessing,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (*SellerWarehouse) TableName() string {
	return "seller_warehouses"
}
