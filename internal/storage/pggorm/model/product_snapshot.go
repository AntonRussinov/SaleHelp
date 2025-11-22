package model

import (
	"salehelp/domain"
	"time"
)

type ProductSnapshot struct {
	ID             int       `gorm:"primaryKey;column:id"`
	ProductID      int       `gorm:"column:product_id;not null"`
	SnapshotTime   time.Time `gorm:"column:snapshot_time;autoCreateTime"`
	Price          float64   `gorm:"column:price"`
	OldPrice       float64   `gorm:"column:old_price"`
	Discount       float64   `gorm:"column:discount"`
	StockTotal     int       `gorm:"column:stock_total"`
	Rating         float64   `gorm:"column:rating"`
	ReviewsCount   int       `gorm:"column:reviews_count"`
	Sales24h       int       `gorm:"column:sales_24h"`
	Sales7d        int       `gorm:"column:sales_7d"`
	PositionSearch int       `gorm:"column:position_search"`
	Extra          []byte    `gorm:"column:extra"` // JSONB
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func ProductSnapshotColumns() []string {
	return []string{
		"id",
		"product_id",
		"snapshot_time",
		"price",
		"old_price",
		"discount",
		"stock_total",
		"rating",
		"reviews_count",
		"sales_24h",
		"sales_7d",
		"position_search",
		"extra",
		"created_at",
	}
}

func (ps *ProductSnapshot) ToDomain() *domain.ProductSnapshot {
	if ps == nil {
		return nil
	}

	return &domain.ProductSnapshot{
		ID:             ps.ID,
		ProductID:      ps.ProductID,
		SnapshotTime:   ps.SnapshotTime,
		Price:          ps.Price,
		OldPrice:       ps.OldPrice,
		Discount:       ps.Discount,
		StockTotal:     ps.StockTotal,
		Rating:         ps.Rating,
		ReviewsCount:   ps.ReviewsCount,
		Sales24h:       ps.Sales24h,
		Sales7d:        ps.Sales7d,
		PositionSearch: ps.PositionSearch,
		Extra:          ps.Extra,
		CreatedAt:      ps.CreatedAt,
	}
}

func (*ProductSnapshot) TableName() string {
	return "product_snapshots"
}
