package cronjob

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
	r.Log = log.With(sl.Module("CronJobRepository"))
	r.DB = db
	return r
}

func toModel(c *domain.CronJob) *model.CronJob {
	if c == nil {
		return nil
	}

	return &model.CronJob{
		ID:           c.ID,
		Name:         c.Name,
		Status:       c.Status,
		ErrorMessage: c.ErrorMessage,
		StartedAt:    c.StartedAt,
		FinishedAt:   c.FinishedAt,
	}
}
