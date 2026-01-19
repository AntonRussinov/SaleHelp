package alerttriggered

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) IncrementAttempts(ctx context.Context, id int) error {
	if err := r.DB.WithContext(ctx).
		Model(&model.AlertTriggered{}).
		Where("id = ?", id).
		UpdateColumn("send_attempts", gorm.Expr("send_attempts + 1")).Error; err != nil {
		return fmt.Errorf("alerts_triggered inc attempts: %w", err)
	}
	return nil
}

func (r *Repository) MarkSent(ctx context.Context, id int, sentAt time.Time) error {
	if err := r.DB.WithContext(ctx).
		Model(&model.AlertTriggeredDelivery{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      "sent",
			"sent_at":     sentAt,
			"last_error":  nil,
			"lease_until": nil,
			"leased_by":   nil,
			"updated_at":  time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("deliveries mark sent (id=%d): %w", id, err)
	}
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id int, errText string) error {
	const maxLen = 2000
	if len(errText) > maxLen {
		errText = errText[:maxLen]
	}

	if err := r.DB.WithContext(ctx).
		Model(&model.AlertTriggeredDelivery{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      "new", // <-- важный выбор: оставляем new, чтобы ретраилось дальше
			"last_error":  errText,
			"lease_until": nil,
			"leased_by":   nil,
			"updated_at":  time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("deliveries mark failed (id=%d): %w", id, err)
	}
	return nil
}

func (r *Repository) LeaseUnsent(
	ctx context.Context,
	workerID string,
	limit int,
	leaseSeconds int,
) ([]*domain.AlertTriggeredDelivery, error) {

	if limit <= 0 {
		limit = 200
	}
	if leaseSeconds <= 0 {
		leaseSeconds = 300
	}

	q := `
		WITH picked AS (
		  SELECT id
		  FROM alerts_triggered_deliveries
		  WHERE status = 'new'
		    AND (lease_until IS NULL OR lease_until < now())
		  ORDER BY created_at ASC
		  FOR UPDATE SKIP LOCKED
		  LIMIT ?
		)
		UPDATE alerts_triggered_deliveries d
		SET
		  lease_until = now() + (? * interval '1 second'),
		  leased_by   = ?,
		  attempts    = d.attempts + 1,
		  updated_at  = now()
		FROM picked
		WHERE d.id = picked.id
		RETURNING d.*;
	`

	var rows []model.AlertTriggeredDelivery
	if err := r.DB.WithContext(ctx).
		Raw(q, limit, leaseSeconds, workerID).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("deliveries lease unsent: %w", err)
	}

	out := make([]*domain.AlertTriggeredDelivery, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}

func (r *Repository) LeaseUnsentByShop(
	ctx context.Context,
	shopID int,
	workerID string,
	limit int,
	leaseSeconds int,
) ([]*domain.AlertTriggered, error) {

	if limit <= 0 {
		limit = 200
	}
	if leaseSeconds <= 0 {
		leaseSeconds = 300
	}

	q := `
		WITH picked AS (
			SELECT id
			FROM alerts_triggered
			WHERE sent_at IS NULL
			  AND shop_id = ?
			  AND (lease_until IS NULL OR lease_until < now())
			ORDER BY created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT ?
		)
		UPDATE alerts_triggered t
		SET
			lease_until   = now() + (? || ' seconds')::interval,
			leased_by     = ?,
			send_attempts = t.send_attempts + 1
		FROM picked
		WHERE t.id = picked.id
		RETURNING t.*;
	`

	var rows []model.AlertTriggered
	if err := r.DB.WithContext(ctx).
		Raw(q, shopID, limit, leaseSeconds, workerID).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("alerts_triggered lease unsent by shop: %w", err)
	}

	out := make([]*domain.AlertTriggered, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
