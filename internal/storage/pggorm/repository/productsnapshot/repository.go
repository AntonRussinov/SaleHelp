package productsnapshot

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
	r.Log = log.With(sl.Module("ProductSnapshotRepository"))
	r.DB = db
	return r
}

// domain → model
func toModel(d *domain.ProductSnapshot) *model.ProductSnapshot {
	if d == nil {
		return nil
	}

	return &model.ProductSnapshot{
		ID:                d.ID,
		ProductID:         d.ProductID,
		SnapshotTime:      d.SnapshotTime,
		CurrencyCode:      d.CurrencyCode,
		BasePrice:         d.BasePrice,
		FinalPrice:        d.FinalPrice,
		ClubPrice:         d.ClubPrice,
		Discount:          d.Discount,
		ClubDiscount:      d.ClubDiscount,
		OldPrice:          d.OldPrice,
		StockTotal:        d.StockTotal,
		StockInWay:        d.StockInWay,
		Rating:            d.Rating,
		ReviewsCount:      d.ReviewsCount,
		Sales24h:          d.Sales24h,
		Sales7d:           d.Sales7d,
		PositionSearch:    d.PositionSearch,
		IsBadTurnover:     d.IsBadTurnover,
		EditableSizePrice: d.EditableSizePrice,
		Extra:             d.Extra,
		CreatedAt:         d.CreatedAt,
	}
}
