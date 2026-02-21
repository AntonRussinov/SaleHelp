package abcsnapshot

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// CreateBatch вставляет снэпшот ABC (только INSERT, immutable).
func (r *Repository) CreateBatch(ctx context.Context, rows []domain.ABCSnapshotRow) error {
	if len(rows) == 0 {
		return nil
	}
	models := make([]*model.ABCSnapshotRow, 0, len(rows))
	for _, row := range rows {
		cat := string(row.Category)
		if cat == "" {
			cat = "C"
		}
		models = append(models, &model.ABCSnapshotRow{
			ShopID:      row.ShopID,
			ProductID:   row.ProductID,
			NmID:        row.NmID,
			Category:    cat,
			Revenue:     row.Revenue,
			SalesCount:  row.SalesCount,
			PeriodDays:  row.PeriodDays,
			PeriodStart: row.PeriodStart,
			PeriodEnd:   row.PeriodEnd,
			CreatedAt:   row.CreatedAt,
		})
	}
	if err := r.DB.WithContext(ctx).Create(&models).Error; err != nil {
		return fmt.Errorf("abc_snapshots create batch: %w", err)
	}
	return nil
}
