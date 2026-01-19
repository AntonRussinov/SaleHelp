package productsizepricesnapshot

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
	r.Log = log.With(sl.Module("ProductSizePriceSnapshotRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.ProductSizePriceSnapshot) *model.ProductSizePriceSnapshot {
	return &model.ProductSizePriceSnapshot{
		ID:                d.ID,
		ProductSnapshotID: d.ProductSnapshotID,
		ProductSizeID:     d.ProductSizeID,
		BasePrice:         d.BasePrice,
		FinalPrice:        d.FinalPrice,
		ClubPrice:         d.ClubPrice,
		Discount:          d.Discount,
		ClubDiscount:      d.ClubDiscount,
		Extra:             d.Extra,
	}
}
