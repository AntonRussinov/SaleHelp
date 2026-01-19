package model

import (
	"salehelp/domain"
	"time"
)

type ProductSnapshot struct {
	ID           int       `gorm:"primaryKey;column:id"`
	ProductID    int       `gorm:"column:product_id;not null"`
	SnapshotTime time.Time `gorm:"column:snapshot_time;not null;autoCreateTime"`

	CurrencyCode *string `gorm:"column:currency_code;type:char(3)"`

	BasePrice    *float64 `gorm:"column:base_price"`
	FinalPrice   *float64 `gorm:"column:final_price"`
	ClubPrice    *float64 `gorm:"column:club_price"`
	Discount     *float64 `gorm:"column:discount"`
	ClubDiscount *float64 `gorm:"column:club_discount"`
	OldPrice     *float64 `gorm:"column:old_price"`

	StockTotal *int `gorm:"column:stock_total"`
	StockInWay *int `gorm:"column:stock_in_way"`

	Rating       *float64 `gorm:"column:rating"`
	ReviewsCount *int     `gorm:"column:reviews_count"`

	Sales24h *int `gorm:"column:sales_24h"`
	Sales7d  *int `gorm:"column:sales_7d"`

	PositionSearch *int `gorm:"column:position_search"`

	IsBadTurnover     *bool `gorm:"column:is_bad_turnover"`
	EditableSizePrice *bool `gorm:"column:editable_size_price"`

	Extra     []byte    `gorm:"column:extra"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func ProductSnapshotColumns() []string {
	return []string{
		"id",
		"product_id",
		"snapshot_time",
		"currency_code",
		"base_price",
		"final_price",
		"club_price",
		"discount",
		"club_discount",
		"old_price",
		"stock_total",
		"stock_in_way",
		"rating",
		"reviews_count",
		"sales_24h",
		"sales_7d",
		"position_search",
		"is_bad_turnover",
		"editable_size_price",
		"extra",
		"created_at",
	}
}

func (m *ProductSnapshot) ToDomain() *domain.ProductSnapshot {
	if m == nil {
		return nil
	}

	return &domain.ProductSnapshot{
		ID:                m.ID,
		ProductID:         m.ProductID,
		SnapshotTime:      m.SnapshotTime,
		CurrencyCode:      m.CurrencyCode,
		BasePrice:         m.BasePrice,
		FinalPrice:        m.FinalPrice,
		ClubPrice:         m.ClubPrice,
		Discount:          m.Discount,
		ClubDiscount:      m.ClubDiscount,
		OldPrice:          m.OldPrice,
		StockTotal:        m.StockTotal,
		StockInWay:        m.StockInWay,
		Rating:            m.Rating,
		ReviewsCount:      m.ReviewsCount,
		Sales24h:          m.Sales24h,
		Sales7d:           m.Sales7d,
		PositionSearch:    m.PositionSearch,
		IsBadTurnover:     m.IsBadTurnover,
		EditableSizePrice: m.EditableSizePrice,
		Extra:             m.Extra,
		CreatedAt:         m.CreatedAt,
	}
}

func (*ProductSnapshot) TableName() string {
	return "product_snapshots"
}
