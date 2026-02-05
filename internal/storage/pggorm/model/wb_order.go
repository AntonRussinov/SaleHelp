package model

import "time"

type WBOrder struct {
	ID              int        `gorm:"primaryKey;column:id"`
	ShopID          int        `gorm:"column:shop_id;not null"`
	Srid            string     `gorm:"column:srid;type:text;not null"`
	NmID            int64      `gorm:"column:nm_id;not null"`
	Barcode         string     `gorm:"column:barcode;type:text"`
	SupplierArticle string     `gorm:"column:supplier_article;type:text"`
	Date            time.Time  `gorm:"column:date;not null"`
	LastChangeDate  time.Time  `gorm:"column:last_change_date;not null"`
	IsCancel        bool       `gorm:"column:is_cancel;not null"`
	CancelDate      *time.Time `gorm:"column:cancel_date"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null"`
}

func (*WBOrder) TableName() string { return "wb_orders" }
