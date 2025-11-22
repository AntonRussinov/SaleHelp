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
func toModel(s *domain.ProductSnapshot) *model.ProductSnapshot {
	return &model.ProductSnapshot{
		ID:             s.ID,
		ProductID:      s.ProductID,
		SnapshotTime:   s.SnapshotTime,
		Price:          s.Price,
		OldPrice:       s.OldPrice,
		Discount:       s.Discount,
		StockTotal:     s.StockTotal,
		Rating:         s.Rating,
		ReviewsCount:   s.ReviewsCount,
		Sales24h:       s.Sales24h,
		Sales7d:        s.Sales7d,
		PositionSearch: s.PositionSearch,
		Extra:          s.Extra,
		CreatedAt:      s.CreatedAt,
	}
}
