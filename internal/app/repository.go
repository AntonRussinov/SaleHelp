package app

import (
	"log/slog"

	"gorm.io/gorm"
)

type Repository struct {
	//Application *application.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	return &Repository{
		//Application:            application.NewRepository(log, db),
	}
}
