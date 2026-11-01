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
		ID:           s.ID,
		UserID:       s.UserID,
		Tier:         s.Tier,
		Status:       s.Status,
		PeriodMonths: s.PeriodMonths,
		StartedAt:    s.StartedAt,
		ExpiresAt:    s.ExpiresAt,
		TrialEndsAt:  s.TrialEndsAt,
		AIQuota:      s.AIQuota,
		AIUsed:       s.AIUsed,
		AIResetsAt:   s.AIResetsAt,
		MaxShops:     s.MaxShops,
		ExtraShops:   s.ExtraShops,
		PriceMinor:   s.PriceMinor,
	}
}

func (r *Repository) WithDB(db *gorm.DB) *Repository {
	cp := *r
	cp.DB = db
	return &cp
}
