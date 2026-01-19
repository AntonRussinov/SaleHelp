package productsizesnapshot

import (
	"log/slog"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	r := &Repository{}
	r.Log = log.With(sl.Module("ProductSizeSnapshotRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.ProductSizeSnapshot) *model.ProductSizeSnapshot {
	return &model.ProductSizeSnapshot{
		ID:                d.ID,
		ProductSnapshotID: d.ProductSnapshotID,
		ProductSizeID:     d.ProductSizeID,
		SKU:               d.SKU,
		Stock:             d.Stock,
		StockByWarehouse:  d.StockByWarehouse,
		Extra:             d.Extra,
	}
}
