package orderrefusal

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
	r.Log = log.With(sl.Module("OrderRefusalRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.OrderRefusal) *model.OrderRefusal {
	if d == nil {
		return nil
	}
	return &model.OrderRefusal{
		ID:              d.ID,
		ShopID:          d.ShopID,
		Srid:            d.Srid,
		NmID:            d.NmID,
		Barcode:         d.Barcode,
		SupplierArticle: d.SupplierArticle,
		Date:            d.Date,
		CancelDate:      d.CancelDate,
		LastChangeDate:  d.LastChangeDate,
		CreatedAt:       d.CreatedAt,
		NotifiedAt:      d.NotifiedAt,
	}
}
