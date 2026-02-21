package wbfeedback

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
	r.Log = log.With(sl.Module("WBFeedbackRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.WBFeedback) *model.WBFeedback {
	if d == nil {
		return nil
	}
	return &model.WBFeedback{
		ID:               d.ID,
		ShopID:           d.ShopID,
		WBFeedbackID:     d.WBFeedbackID,
		NmID:             d.NmID,
		ProductName:      d.ProductName,
		Text:             d.Text,
		Pros:             d.Pros,
		Cons:             d.Cons,
		ProductValuation: d.ProductValuation,
		UserName:         d.UserName,
		CreatedAt:        d.CreatedAt,
		NotifiedAt:       d.NotifiedAt,
		IsReplied:        d.IsReplied,
		RepliedAt:        d.RepliedAt,
		OrderStatus:      d.OrderStatus,
		Extra:            d.Extra,
	}
}
