package ai_feedback_reply

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
	r.Log = log.With(sl.Module("AIFeedbackReplyRepository"))
	r.DB = db
	return r
}
