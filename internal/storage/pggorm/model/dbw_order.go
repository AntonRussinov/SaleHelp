package model

import (
	"salehelp/domain"
	"time"
)

type DBWOrder struct {
	ID             int        `gorm:"primaryKey;column:id"`
	ShopID         int        `gorm:"column:shop_id;not null"`
	WBOrderID      int64      `gorm:"column:wb_order_id;not null"`
	OrderUID       string     `gorm:"column:order_uid;type:text;not null"`
	NmID           int64      `gorm:"column:nm_id;not null"`
	ChrtID         int64      `gorm:"column:chrt_id;not null"`
	Article        string     `gorm:"column:article;type:text"`
	SupplierStatus string     `gorm:"column:supplier_status;type:text"`
	WBStatus       string     `gorm:"column:wb_status;type:text"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null"`
	NotifiedAt     *time.Time `gorm:"column:notified_at"`
	Extra          []byte     `gorm:"column:extra;type:jsonb"`
}

func (m *DBWOrder) ToDomain() *domain.DBWOrder {
	if m == nil {
		return nil
	}
	return &domain.DBWOrder{
		ID:             m.ID,
		ShopID:         m.ShopID,
		WBOrderID:      m.WBOrderID,
		OrderUID:       m.OrderUID,
		NmID:           m.NmID,
		ChrtID:         m.ChrtID,
		Article:        m.Article,
		SupplierStatus: m.SupplierStatus,
		WBStatus:       m.WBStatus,
		CreatedAt:      m.CreatedAt,
		NotifiedAt:     m.NotifiedAt,
		Extra:          m.Extra,
	}
}

func (*DBWOrder) TableName() string {
	return "dbw_orders"
}
