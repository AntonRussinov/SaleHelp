package saleevent

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
	r.Log = log.With(sl.Module("SaleEventRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.SaleEvent) *model.SaleEvent {
	if d == nil {
		return nil
	}
	eventType := string(d.EventType)
	if eventType == "" {
		eventType = string(domain.SaleEventTypeBuyout)
	}
	return &model.SaleEvent{
		ID:              d.ID,
		ShopID:          d.ShopID,
		EventType:       eventType,
		WBSaleID:        d.WBSaleID,
		Srid:            d.Srid,
		Date:            d.Date,
		LastChangeDate:  d.LastChangeDate,
		NmID:            d.NmID,
		Barcode:         d.Barcode,
		SupplierArticle: d.SupplierArticle,
		WarehouseName:   d.WarehouseName,
		Category:        d.Category,
		Subject:         d.Subject,
		Brand:           d.Brand,
		TechSize:        d.TechSize,
		IncomeID:        d.IncomeID,
		TotalPrice:      d.TotalPrice,
		DiscountPercent: d.DiscountPercent,
		FinishedPrice:   d.FinishedPrice,
		PriceWithDisc:   d.PriceWithDisc,
		ForPay:          d.ForPay,
		CreatedAt:       d.CreatedAt,
		NotifiedAt:      d.NotifiedAt,
		Extra:           d.Extra,
	}
}
