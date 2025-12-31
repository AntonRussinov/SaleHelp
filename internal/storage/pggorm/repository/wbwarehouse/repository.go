package wbwarehouse

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
	r.Log = log.With(sl.Module("WBWarehouseRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.WBWarehouse) *model.WBWarehouse {
	return &model.WBWarehouse{
		ID:              d.ID,
		WBID:            d.WBID,
		Name:            d.Name,
		Address:         d.Address,
		City:            d.City,
		Longitude:       d.Longitude,
		Latitude:        d.Latitude,
		CargoType:       d.CargoType,
		DeliveryType:    d.DeliveryType,
		FederalDistrict: d.FederalDistrict,
		Selected:        d.Selected,
	}
}
