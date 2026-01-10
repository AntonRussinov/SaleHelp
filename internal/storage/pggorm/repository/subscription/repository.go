package subscription

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
	repo.Log = log.With(sl.Module("SubscriptionRepository"))
	repo.DB = db
	return repo
}

func toModel(s *domain.Subscription) *model.Subscription {
	if s == nil {
		return nil
	}
	return &model.Subscription{
		ID:              s.ID,
		ShopID:          s.ShopID,
		Plan:            s.Plan,
		Status:          s.Status,
		ValidFrom:       s.ValidFrom,
		ValidTo:         s.ValidTo,
		PeriodMonths:    s.PeriodMonths,
		BaseAmountMinor: s.BaseAmountMinor,
		CurrencyCode:    s.CurrencyCode,
		PaidAt:          s.PaidAt,
		Source:          s.Source,
		Comment:         s.Comment,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}
