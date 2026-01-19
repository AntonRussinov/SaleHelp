package product

import (
	"log/slog"

	"salehelp/domain"

	"salehelp/internal/storage/pggorm" // твой базовый репозиторий
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	repo := &Repository{}
	repo.Log = log.With(sl.Module("ProductRepository"))
	repo.DB = db

	return repo
}

// helper: domain → model
func toModel(p *domain.Product) *model.Product {
	if p == nil {
		return nil
	}

	return &model.Product{
		ID:          p.ID,
		ShopID:      p.ShopID,
		ExternalID:  p.ExternalID,
		Title:       p.Title,
		URL:         p.URL,
		ImageURL:    p.ImageURL,
		Brand:       p.Brand,
		Category:    p.Category,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
