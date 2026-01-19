package notificationendpoint

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
	r.Log = log.With(sl.Module("NotificationEndpointRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.NotificationEndpoint) *model.ShopNotificationEndpoint {
	if d == nil {
		return nil
	}

	return &model.ShopNotificationEndpoint{
		ID:      d.ID,
		ShopID:  d.ShopID,
		Channel: d.Channel,
		Address: d.Address,
		Enabled: d.Enabled,
		Label:   d.Label,
	}
}
