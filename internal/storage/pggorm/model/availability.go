package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

type Availability struct {
	ID           int            `gorm:"primaryKey;column:id"`
	MerchantUID  string         `gorm:"column:merchant_uid;type:varchar(100)"`
	OfferID      int            `gorm:"column:offer_id"`
	Available    bool           `gorm:"column:available"`
	StoreID      string         `gorm:"column:store_id"`
	PreOrder     int            `gorm:"column:pre_order"`
	StockCount   int            `gorm:"column:stock_count"`
	StockEnabled bool           `gorm:"column:stock_enabled"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func AvailabilityColumns() []string {
	return []string{
		"id",
		"merchant_uid",
		"offer_id",
		"available",
		"store_id",
		"pre_order",
		"stock_count",
		"stock_enabled",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (a *Availability) ToDomain() *domain.Availability {
	if a == nil {
		return nil
	}
	return &domain.Availability{
		ID:           a.ID,
		MerchantUID:  a.MerchantUID,
		OfferID:      a.OfferID,
		Available:    a.Available,
		StoreID:      a.StoreID,
		PreOrder:     a.PreOrder,
		StockCount:   a.StockCount,
		StockEnabled: a.StockEnabled,
	}
}

func (*Availability) TableName() string {
	return "availabilities"
}
