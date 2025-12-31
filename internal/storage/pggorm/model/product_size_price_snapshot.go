package model

import "salehelp/domain"

type ProductSizePriceSnapshot struct {
	ID                int `gorm:"primaryKey;column:id"`
	ProductSnapshotID int `gorm:"column:product_snapshot_id;not null"`

	ProductSizeID int `gorm:"column:product_size_id;not null"`

	BasePrice    *float64 `gorm:"column:base_price"`
	FinalPrice   *float64 `gorm:"column:final_price"`
	ClubPrice    *float64 `gorm:"column:club_price"`
	Discount     *float64 `gorm:"column:discount"`
	ClubDiscount *float64 `gorm:"column:club_discount"`

	Extra []byte `gorm:"column:extra"`
}

func ProductSizePriceSnapshotColumns() []string {
	return []string{
		"id",
		"product_snapshot_id",
		"product_size_id",
		"base_price",
		"final_price",
		"club_price",
		"discount",
		"club_discount",
		"extra",
	}
}

func (p *ProductSizePriceSnapshot) ToDomain() *domain.ProductSizePriceSnapshot {
	if p == nil {
		return nil
	}
	return &domain.ProductSizePriceSnapshot{
		ID:                p.ID,
		ProductSnapshotID: p.ProductSnapshotID,
		ProductSizeID:     p.ProductSizeID,
		BasePrice:         p.BasePrice,
		FinalPrice:        p.FinalPrice,
		ClubPrice:         p.ClubPrice,
		Discount:          p.Discount,
		ClubDiscount:      p.ClubDiscount,
		Extra:             p.Extra,
	}
}

func (*ProductSizePriceSnapshot) TableName() string {
	return "product_size_price_snapshots"
}
