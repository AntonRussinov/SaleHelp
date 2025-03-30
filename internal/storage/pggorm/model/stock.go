package model

import (
	"salehelp/domain"
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	OfferID     int            `gorm:"column:offer_id;primaryKey;not null"`            // Ссылка на предложение из таблицы offers
	MerchantUID string         `gorm:"column:merchant_uid;type:varchar(100);not null"` // Идентификатор продавца из таблицы merchants
	StoreID     string         `gorm:"column:store_id;primaryKey;not null"`            // Внешний идентификатор склада/точки (store_id)
	CityID      string         `gorm:"column:city_id"`                                 // Идентификатор города (если применимо)
	Available   bool           `gorm:"column:available"`                               // Флаг доступности (например, преобразованный из "yes")
	PreOrder    int            `gorm:"column:pre_order"`                               // Количество товара для предзаказа (если значение числовое)
	StockCount  int            `gorm:"column:stock_count"`                             // Количество товара на складе
	Enabled     bool           `gorm:"column:enabled"`                                 // Флаг, указывающий, активна ли запись для данного склада
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`               // Дата создания записи
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`               // Дата последнего обновления
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`                        // Мягкое удаление (soft delete)
}

func StockColumns() []string {
	return []string{
		"merchant_uid",
		"offer_id",
		"city_id",
		"store_id",
		"available",
		"pre_order",
		"stock_count",
		"enabled",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (s *Stock) ToDomain() *domain.Stock {
	if s == nil {
		return nil
	}
	return &domain.Stock{
		MerchantUID: s.MerchantUID,
		OfferID:     s.OfferID,
		CityID:      s.CityID,
		StoreID:     s.StoreID,
		StockCount:  s.StockCount,
		PreOrder:    s.PreOrder,
		Available:   s.Available,
		Enabled:     s.Enabled,
	}
}

func (*Stock) TableName() string {
	return "stocks"
}
