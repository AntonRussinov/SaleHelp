package shopnotificationsetting

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
	repo := &Repository{}
	repo.Log = log.With(sl.Module("ShopNotificationSettingRepository"))
	repo.DB = db
	return repo
}

func (r *Repository) WithDB(db *gorm.DB) *Repository {
	cp := *r
	cp.DB = db
	return &cp
}
