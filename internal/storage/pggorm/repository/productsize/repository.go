package productsize

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
	r.Log = log.With(sl.Module("ProductSizeRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.ProductSize) *model.ProductSize {
	return &model.ProductSize{
		ID:           d.ID,
		ProductID:    d.ProductID,
		ExternalID:   d.ExternalID,
		TechSizeName: d.TechSizeName,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
