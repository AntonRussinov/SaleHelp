package model

import (
	"salehelp/domain"
	"time"
)

type WBWarehouse struct {
	ID              int       `gorm:"primaryKey;column:id"`
	WBID            int64     `gorm:"column:wb_id;not null;unique"`
	Name            string    `gorm:"column:name;not null"`
	Address         *string   `gorm:"column:address"`
	City            *string   `gorm:"column:city"`
	Longitude       *float64  `gorm:"column:longitude"`
	Latitude        *float64  `gorm:"column:latitude"`
	CargoType       *int      `gorm:"column:cargo_type"`
	DeliveryType    *int      `gorm:"column:delivery_type"`
	FederalDistrict *string   `gorm:"column:federal_district"`
	Selected        *bool     `gorm:"column:selected"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func WBWarehouseColumns() []string {
	return []string{
		"id",
		"wb_id",
		"name",
		"address",
		"city",
		"longitude",
		"latitude",
		"cargo_type",
		"delivery_type",
		"federal_district",
		"selected",
		"created_at",
		"updated_at",
	}
}

func (m *WBWarehouse) ToDomain() *domain.WBWarehouse {
	return &domain.WBWarehouse{
		ID:              m.ID,
		WBID:            m.WBID,
		Name:            m.Name,
		Address:         m.Address,
		City:            m.City,
		Longitude:       m.Longitude,
		Latitude:        m.Latitude,
		CargoType:       m.CargoType,
		DeliveryType:    m.DeliveryType,
		FederalDistrict: m.FederalDistrict,
		Selected:        m.Selected,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (*WBWarehouse) TableName() string {
	return "wb_warehouses"
}
