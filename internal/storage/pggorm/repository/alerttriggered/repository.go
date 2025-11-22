package alerttriggered

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
	r.Log = log.With(sl.Module("AlertTriggeredRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.AlertTriggered) *model.AlertTriggered {
	if d == nil {
		return nil
	}

	return &model.AlertTriggered{
		ID:         d.ID,
		AlertID:    d.AlertID,
		ProductID:  d.ProductID,
		SnapshotID: d.SnapshotID,
		Message:    d.Message,
		CreatedAt:  d.CreatedAt,
	}
}
