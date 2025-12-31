package sellerwarehouse

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
	r.Log = log.With(sl.Module("SellerWarehouseRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.SellerWarehouse) *model.SellerWarehouse {
	return &model.SellerWarehouse{
		ID:           d.ID,
		ShopID:       d.ShopID,
		WBOfficeID:   d.WBOfficeID,
		WBID:         d.WBID,
		Name:         d.Name,
		CargoType:    d.CargoType,
		DeliveryType: d.DeliveryType,
		IsDeleting:   d.IsDeleting,
		IsProcessing: d.IsProcessing,
	}
}
