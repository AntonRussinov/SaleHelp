package abcsnapshot

import (
	"log/slog"

	"salehelp/internal/storage/pggorm"
	"salehelp/pkg/logger/sl"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	r := &Repository{}
	r.Log = log.With(sl.Module("ABCSnapshotRepository"))
	r.DB = db
	return r
}
