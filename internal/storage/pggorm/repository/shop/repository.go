package shop

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
	repo := &Repository{}
	repo.Log = log.With(sl.Module("ShopRepository"))
	repo.DB = db

	return repo
}

// domain → model
func toModel(s *domain.Shop) *model.Shop {
	if s == nil {
		return nil
	}

	return &model.Shop{
		ID:          s.ID,
		UserID:      s.UserID,
		Marketplace: s.Marketplace,
		ExternalID:  s.ExternalID,
		Name:        s.Name,
		APIToken:    s.APIToken,
		IsActive:    s.IsActive,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
