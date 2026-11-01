package model

import (
	"salehelp/domain"
	"time"
)

type SaleEvent struct {
	ID              int        `gorm:"primaryKey;column:id"`
	ShopID          int        `gorm:"column:shop_id;not null"`
	EventType       string     `gorm:"column:event_type;type:sale_event_type;not null"`
	WBSaleID        string     `gorm:"column:wb_sale_id;type:text;not null"`
	Srid            string     `gorm:"column:srid;type:text;not null"`
	Date            time.Time  `gorm:"column:date;not null"`
	LastChangeDate  time.Time  `gorm:"column:last_change_date;not null"`
	NmID            int64      `gorm:"column:nm_id;not null"`
	Barcode         string     `gorm:"column:barcode;type:text"`
	SupplierArticle string     `gorm:"column:supplier_article;type:text"`
	WarehouseName   string     `gorm:"column:warehouse_name;type:text"`
	Category        string     `gorm:"column:category;type:text"`
	Subject         string     `gorm:"column:subject;type:text"`
	Brand           string     `gorm:"column:brand;type:text"`
	TechSize        string     `gorm:"column:tech_size;type:text"`
	IncomeID        int64      `gorm:"column:income_id"`
	TotalPrice      float64    `gorm:"column:total_price"`
	DiscountPercent float64    `gorm:"column:discount_percent"`
	FinishedPrice   float64    `gorm:"column:finished_price"`
	PriceWithDisc   float64    `gorm:"column:price_with_disc"`
	ForPay          float64    `gorm:"column:for_pay"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null"`
	NotifiedAt      *time.Time `gorm:"column:notified_at"`
	Extra           []byte     `gorm:"column:extra;type:jsonb"`
}

func (m *SaleEvent) ToDomain() *domain.SaleEvent {
	if m == nil {
		return nil
	}
	return &domain.SaleEvent{
		ID:              m.ID,
		ShopID:          m.ShopID,
		EventType:       domain.SaleEventType(m.EventType),
		WBSaleID:        m.WBSaleID,
		Srid:            m.Srid,
		Date:            m.Date,
		LastChangeDate:  m.LastChangeDate,
		NmID:            m.NmID,
		Barcode:         m.Barcode,
		SupplierArticle: m.SupplierArticle,
		WarehouseName:   m.WarehouseName,
		Category:        m.Category,
		Subject:         m.Subject,
		Brand:           m.Brand,
		TechSize:        m.TechSize,
		IncomeID:        m.IncomeID,
		TotalPrice:      m.TotalPrice,
		DiscountPercent: m.DiscountPercent,
		FinishedPrice:   m.FinishedPrice,
		PriceWithDisc:   m.PriceWithDisc,
		ForPay:          m.ForPay,
		CreatedAt:       m.CreatedAt,
		NotifiedAt:      m.NotifiedAt,
		Extra:           m.Extra,
	}
}

func (*SaleEvent) TableName() string {
	return "sale_events"
}
