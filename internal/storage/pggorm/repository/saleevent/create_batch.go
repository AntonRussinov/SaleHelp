package saleevent

import (
	"context"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

const saleEventsBatchSize = 500

// CreateBatch вставляет записи батчами с ON CONFLICT (shop_id, wb_sale_id) DO NOTHING.
// У одного заказа (srid) бывают и выкуп (saleID S...), и возврат (saleID R...) — уникаль по wb_sale_id.
func (r *Repository) CreateBatch(ctx context.Context, events []*domain.SaleEvent) error {
	if len(events) == 0 {
		return nil
	}
	models := make([]*model.SaleEvent, 0, len(events))
	for _, e := range events {
		models = append(models, toModel(e))
	}
	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shop_id"}, {Name: "wb_sale_id"}},
			DoNothing: true,
		}).
		CreateInBatches(models, saleEventsBatchSize).
		Error
	if err != nil {
		return err
	}
	return nil
}
