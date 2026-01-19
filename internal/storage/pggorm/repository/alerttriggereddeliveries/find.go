package alerttriggereddeliveries

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.AlertTriggeredDelivery, error) {
	var m model.AlertTriggeredDelivery

	if err := r.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("alert_triggered_delivery get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByTriggeredID(ctx context.Context, triggeredID int) ([]*domain.AlertTriggeredDelivery, error) {
	var rows []model.AlertTriggeredDelivery

	if err := r.DB.WithContext(ctx).
		Where("triggered_id = ?", triggeredID).
		Order("id").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("alert_triggered_delivery list by triggered: %w", err)
	}

	out := make([]*domain.AlertTriggeredDelivery, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}

// ListReady — то, что готово к отправке (new/failed) и без активной аренды
func (r *Repository) ListReady(ctx context.Context, limit int) ([]*domain.AlertTriggeredDelivery, error) {
	var rows []model.AlertTriggeredDelivery

	q := r.DB.WithContext(ctx).
		Where(`status IN ('new','failed')`).
		Where(`lease_until IS NULL OR lease_until < now()`).
		Order("id")

	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alert_triggered_delivery list ready: %w", err)
	}

	out := make([]*domain.AlertTriggeredDelivery, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
