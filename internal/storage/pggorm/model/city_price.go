package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

type CityPrice struct {
	ID          int            `gorm:"primaryKey;column:id"`
	MerchantUID string         `gorm:"column:merchant_uid;type:varchar(100)"`
	OfferID     int            `gorm:"column:offer_id"`
	CityID      string         `gorm:"column:city_id"`
	Price       float64        `gorm:"column:price"`
	OldPrice    float64        `gorm:"column:old_price"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func CityPriceColumns() []string {
	return []string{
		"id",
		"merchant_uid",
		"offer_id",
		"city_id",
		"price",
		"old_price",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (cp *CityPrice) ToDomain() *domain.CityPrice {
	if cp == nil {
		return nil
	}
	return &domain.CityPrice{
		ID:          cp.ID,
		MerchantUID: cp.MerchantUID,
		OfferID:     cp.OfferID,
		CityID:      cp.CityID,
		Price:       cp.Price,
		OldPrice:    cp.OldPrice,
	}
}

func (*CityPrice) TableName() string {
	return "city_prices"
}
