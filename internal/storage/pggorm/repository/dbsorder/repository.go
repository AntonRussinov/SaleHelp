package dbsorder

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
	r.Log = log.With(sl.Module("DBSOrderRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.DBSOrder) *model.DBSOrder {
	if d == nil {
		return nil
	}
	return &model.DBSOrder{
		ID:             d.ID,
		ShopID:         d.ShopID,
		WBOrderID:      d.WBOrderID,
		OrderUID:       d.OrderUID,
		NmID:           d.NmID,
		ChrtID:         d.ChrtID,
		Article:        d.Article,
		SupplierStatus: d.SupplierStatus,
		WBStatus:       d.WBStatus,
		CreatedAt:      d.CreatedAt,
		NotifiedAt:     d.NotifiedAt,
		Extra:          d.Extra,
	}
}
