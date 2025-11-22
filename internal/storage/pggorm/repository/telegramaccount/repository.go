package telegramaccount

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
	r.Log = log.With(sl.Module("TelegramAccountRepository"))
	r.DB = db
	return r
}

func toModel(a *domain.TelegramAccount) *model.TelegramAccount {
	if a == nil {
		return nil
	}

	return &model.TelegramAccount{
		ID:         a.ID,
		UserID:     a.UserID,
		TelegramID: a.TelegramID,
		ChatID:     a.ChatID,
		Username:   a.Username,
		CreatedAt:  a.CreatedAt,
	}
}
