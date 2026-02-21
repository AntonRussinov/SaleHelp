package model

import (
	"salehelp/domain"
	"time"
)

type ABCSnapshotRow struct {
	ID         int        `gorm:"primaryKey;column:id"`
	ShopID     int        `gorm:"column:shop_id;not null"`
	ProductID  *int       `gorm:"column:product_id"`
	NmID       int64      `gorm:"column:nm_id;not null"`
	Category   string     `gorm:"column:category;type:char(1);not null"`
	Revenue    float64    `gorm:"column:revenue;not null"`
	SalesCount int        `gorm:"column:sales_count;not null"`
	PeriodDays int        `gorm:"column:period_days;not null"`
	PeriodStart time.Time `gorm:"column:period_start;not null"`
	PeriodEnd   time.Time `gorm:"column:period_end;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
}

func (m *ABCSnapshotRow) ToDomain() *domain.ABCSnapshotRow {
	if m == nil {
		return nil
	}
	return &domain.ABCSnapshotRow{
		ID:         m.ID,
		ShopID:     m.ShopID,
		ProductID:  m.ProductID,
		NmID:       m.NmID,
		Category:   domain.ABCCategory(m.Category),
		Revenue:    m.Revenue,
		SalesCount: m.SalesCount,
		PeriodDays: m.PeriodDays,
		PeriodStart: m.PeriodStart,
		PeriodEnd:   m.PeriodEnd,
		CreatedAt:   m.CreatedAt,
	}
}

func (*ABCSnapshotRow) TableName() string {
	return "abc_snapshots"
}
