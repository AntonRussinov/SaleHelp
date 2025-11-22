package airecommendation

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
	r.Log = log.With(sl.Module("AIRecommendationRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.AIRecommendation) *model.AIRecommendation {
	return &model.AIRecommendation{
		ID:        d.ID,
		ShopID:    d.ShopID,
		ProductID: d.ProductID,
		AlertID:   d.AlertID,
		Type:      d.Type,
		Model:     d.Model,
		Prompt:    d.Prompt,
		Content:   d.Content,
		Meta:      d.Meta,
		CreatedAt: d.CreatedAt,
	}
}
