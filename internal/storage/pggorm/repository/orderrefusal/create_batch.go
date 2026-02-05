package orderrefusal

import (
	"context"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

const batchSize = 500

// CreateBatch вставляет отказы батчами с ON CONFLICT (shop_id, srid) DO NOTHING.
func (r *Repository) CreateBatch(ctx context.Context, refusals []*domain.OrderRefusal) error {
	if len(refusals) == 0 {
		return nil
	}
	models := make([]*model.OrderRefusal, 0, len(refusals))
	for _, e := range refusals {
		models = append(models, toModel(e))
	}
	return r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shop_id"}, {Name: "srid"}},
			DoNothing: true,
		}).
		CreateInBatches(models, batchSize).
		Error
}
