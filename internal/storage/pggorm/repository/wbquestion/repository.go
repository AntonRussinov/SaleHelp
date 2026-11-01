package wbquestion

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
	r.Log = log.With(sl.Module("WBQuestionRepository"))
	r.DB = db
	return r
}

func toModel(d *domain.WBQuestion) *model.WBQuestion {
	if d == nil {
		return nil
	}
	return &model.WBQuestion{
		ID:            d.ID,
		ShopID:        d.ShopID,
		WBQuestionID:  d.WBQuestionID,
		NmID:          d.NmID,
		ProductName:   d.ProductName,
		Text:          d.Text,
		CreatedAt:     d.CreatedAt,
		NotifiedAt:    d.NotifiedAt,
		Extra:         d.Extra,
	}
}
