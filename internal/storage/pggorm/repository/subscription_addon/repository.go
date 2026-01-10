package subscriptionaddon

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
	repo.Log = log.With(sl.Module("SubscriptionAddonRepository"))
	repo.DB = db
	return repo
}

func toModel(a *domain.SubscriptionAddon) *model.SubscriptionAddon {
	if a == nil {
		return nil
	}
	meta := a.Meta
	if len(meta) == 0 {
		meta = []byte(`{}`)
	}
	return &model.SubscriptionAddon{
		ID:             a.ID,
		SubscriptionID: a.SubscriptionID,
		Code:           a.Code,
		Quantity:       a.Quantity,
		AmountMinor:    a.AmountMinor,
		CurrencyCode:   a.CurrencyCode,
		Meta:           meta,
		CreatedAt:      a.CreatedAt,
	}
}
