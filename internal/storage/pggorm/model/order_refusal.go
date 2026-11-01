package model

import (
	"salehelp/domain"
	"time"
)

type OrderRefusal struct {
	ID              int        `gorm:"primaryKey;column:id"`
	ShopID          int        `gorm:"column:shop_id;not null"`
	Srid            string     `gorm:"column:srid;type:text;not null"`
	NmID            int64      `gorm:"column:nm_id;not null"`
	Barcode         string     `gorm:"column:barcode;type:text"`
	SupplierArticle string     `gorm:"column:supplier_article;type:text"`
	Date            time.Time  `gorm:"column:date;not null"`
	CancelDate      time.Time  `gorm:"column:cancel_date;not null"`
	LastChangeDate  time.Time  `gorm:"column:last_change_date;not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null"`
	NotifiedAt      *time.Time `gorm:"column:notified_at"`
}

func (m *OrderRefusal) ToDomain() *domain.OrderRefusal {
	if m == nil {
		return nil
	}
	return &domain.OrderRefusal{
		ID:              m.ID,
		ShopID:          m.ShopID,
		Srid:            m.Srid,
		NmID:            m.NmID,
		Barcode:         m.Barcode,
		SupplierArticle: m.SupplierArticle,
		Date:            m.Date,
		CancelDate:      m.CancelDate,
		LastChangeDate:  m.LastChangeDate,
		CreatedAt:       m.CreatedAt,
		NotifiedAt:      m.NotifiedAt,
	}
}

func (*OrderRefusal) TableName() string {
	return "order_refusals"
}
