package model

import "salehelp/domain"

type ProductSizeSnapshot struct {
	ID                int    `gorm:"primaryKey;column:id"`
	ProductSnapshotID int    `gorm:"column:product_snapshot_id;not null"`
	ProductSizeID     int    `gorm:"column:product_size_id;not null"`
	SKU               string `gorm:"column:sku;not null"`

	Stock            *int   `gorm:"column:stock"`
	StockByWarehouse []byte `gorm:"column:stock_by_warehouse"`
	Extra            []byte `gorm:"column:extra"`
}

func ProductSizeSnapshotColumns() []string {
	return []string{
		"id",
		"product_snapshot_id",
		"product_size_id",
		"sku",
		"stock",
		"stock_by_warehouse",
		"extra",
	}
}

func (p *ProductSizeSnapshot) ToDomain() *domain.ProductSizeSnapshot {
	if p == nil {
		return nil
	}
	return &domain.ProductSizeSnapshot{
		ID:                p.ID,
		ProductSnapshotID: p.ProductSnapshotID,
		ProductSizeID:     p.ProductSizeID,
		SKU:               p.SKU,
		Stock:             p.Stock,
		StockByWarehouse:  p.StockByWarehouse,
		Extra:             p.Extra,
	}
}

func (*ProductSizeSnapshot) TableName() string {
	return "product_size_snapshots"
}
