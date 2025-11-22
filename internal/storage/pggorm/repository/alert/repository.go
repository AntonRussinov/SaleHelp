package alert

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
	r.Log = log.With(sl.Module("AlertRepository"))
	r.DB = db
	return r
}

func toModel(a *domain.Alert) *model.Alert {
	if a == nil {
		return nil
	}

	return &model.Alert{
		ID:                 a.ID,
		ShopID:             a.ShopID,
		ProductID:          a.ProductID,
		Type:               a.Type,
		Severity:           a.Severity,
		Status:             a.Status,
		Message:            a.Message,
		Payload:            a.Payload,
		TriggeredAt:        a.TriggeredAt,
		ResolvedAt:         a.ResolvedAt,
		AIRecommendationID: a.AIRecommendationID,
	}
}
