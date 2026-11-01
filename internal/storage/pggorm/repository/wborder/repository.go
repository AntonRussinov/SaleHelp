package wborder

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"salehelp/internal/storage/pggorm"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	r := &Repository{}
	r.Log = log.With(sl.Module("WBOrderRepository"))
	r.DB = db
	return r
}

func (r *Repository) CountByShop(ctx context.Context, shopID int) (int, error) {
	var cnt int64
	if err := r.DB.WithContext(ctx).
		Model(&model.WBOrder{}).
		Where("shop_id = ?", shopID).
		Count(&cnt).Error; err != nil {
		return 0, err
	}
	return int(cnt), nil
}

// UpsertBatch сохраняет/обновляет заказы по уникальному ключу (shop_id, srid).
func (r *Repository) UpsertBatch(ctx context.Context, rows []*model.WBOrder) error {
	if len(rows) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "shop_id"}, {Name: "srid"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"nm_id",
				"barcode",
				"supplier_article",
				"date",
				"last_change_date",
				"is_cancel",
				"cancel_date",
				"updated_at",
			}),
		}).
		CreateInBatches(rows, 500).
		Error
}

// DailyDemandByNmID возвращает средний дневной "спрос" по nm_id за период since..now:
// +1 за каждый заказ (date), -1 за отмену/отказ (is_cancel=true, cancel_date).
func (r *Repository) DailyDemandByNmID(ctx context.Context, shopID int, since time.Time) (map[int64]float64, error) {
	type row struct {
		NmID     int64
		NetCount float64
	}
	var rows []row

	// Делим на число дней в окне, аналогично sale_events.DailyConsumptionByNmID
	err := r.DB.WithContext(ctx).
		Raw(`
			SELECT nm_id,
				(
					SUM(CASE WHEN date >= ? THEN 1.0 ELSE 0 END)
					-
					SUM(CASE WHEN is_cancel = true AND cancel_date IS NOT NULL AND cancel_date >= ? THEN 1.0 ELSE 0 END)
				) / NULLIF(EXTRACT(EPOCH FROM (NOW() - ?)) / 86400, 0) AS net_count
			FROM wb_orders
			WHERE shop_id = ? AND (date >= ? OR (is_cancel = true AND cancel_date IS NOT NULL AND cancel_date >= ?))
			GROUP BY nm_id
		`, since, since, since, shopID, since, since).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("wb_orders daily demand by nm_id: %w", err)
	}
	out := make(map[int64]float64, len(rows))
	for _, rr := range rows {
		if rr.NetCount > 0 {
			out[rr.NmID] = rr.NetCount
		}
	}
	return out, nil
}
