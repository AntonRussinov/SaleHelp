package productsizesku

import (
	"log/slog"

	"salehelp/internal/storage/pggorm"
	"salehelp/pkg/logger/sl"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	r := &Repository{}
	r.Log = log.With(sl.Module("ProductSizeSKURepository"))
	r.DB = db
	return r
}

// func toModel(d *domain.ProductSizeSKU) *model.ProductSizeSKU {
// 	return &model.ProductSizeSKU{
// 		ID:            d.ID,
// 		ProductSizeID: d.ProductSizeID,
// 		SKU:           d.SKU,
// 		CreatedAt:     d.CreatedAt,
// 	}
// }
