package alerttriggered

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.AlertTriggered, error) {
	var m model.AlertTriggered

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("alert_triggered get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByAlertID(ctx context.Context, alertID int, limit, offset int) ([]*domain.AlertTriggered, error) {
	var rows []model.AlertTriggered

	q := r.DB.WithContext(ctx).
		Where("alert_id = ?", alertID).
		Order("created_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alert_triggered list by alert: %w", err)
	}

	out := make([]*domain.AlertTriggered, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) ListByProductID(ctx context.Context, productID int, limit, offset int) ([]*domain.AlertTriggered, error) {
	var rows []model.AlertTriggered

	q := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alert_triggered list by product: %w", err)
	}

	out := make([]*domain.AlertTriggered, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) ListUnsent(ctx context.Context, limit int) ([]*domain.AlertTriggered, error) {
	if limit <= 0 {
		limit = 100
	}

	var rows []model.AlertTriggered
	if err := r.DB.WithContext(ctx).
		Where("sent_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("alerts_triggered list unsent: %w", err)
	}

	out := make([]*domain.AlertTriggered, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
