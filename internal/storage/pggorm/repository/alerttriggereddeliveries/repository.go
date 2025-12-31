package alerttriggereddeliveries

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
	r.Log = log.With(sl.Module("AlertTriggeredDeliveryRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.AlertTriggeredDelivery) *model.AlertTriggeredDelivery {
	if d == nil {
		return nil
	}

	return &model.AlertTriggeredDelivery{
		ID:          d.ID,
		TriggeredID: d.TriggeredID,
		EndpointID:  d.EndpointID,
		Status:      d.Status,
		SentAt:      d.SentAt,
		Attempts:    d.Attempts,
		LastError:   d.LastError,
		LeaseUntil:  d.LeaseUntil,
		LeasedBy:    d.LeasedBy,
	}
}
