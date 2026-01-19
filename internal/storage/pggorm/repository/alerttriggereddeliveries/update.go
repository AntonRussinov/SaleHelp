package alerttriggereddeliveries

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"github.com/jackc/pgtype"
)

func (r *Repository) Update(ctx context.Context, d *domain.AlertTriggeredDelivery) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("alert_triggered_delivery update: %w", err)
	}

	d.UpdatedAt = m.UpdatedAt
	return nil
}

// LeaseBatch — безопасно "забрать" пачку задач на отправку (одним запросом)
// Это удобно для воркера/крона: один воркер — один lease_owner.
func (r *Repository) LeaseBatch(ctx context.Context, leaseOwner string, leaseFor time.Duration, limit int) ([]*domain.AlertTriggeredDelivery, error) {
	if limit <= 0 {
		limit = 100
	}

	// UPDATE ... FROM (SELECT ... FOR UPDATE SKIP LOCKED) RETURNING *
	leaseUntil := time.Now().Add(leaseFor)

	var leased []model.AlertTriggeredDelivery
	err := r.DB.WithContext(ctx).Raw(`
		WITH cte AS (
			SELECT id
			FROM alerts_triggered_deliveries
			WHERE status IN ('new','failed')
			  AND (lease_until IS NULL OR lease_until < now())
			ORDER BY id
			FOR UPDATE SKIP LOCKED
			LIMIT ?
		)
		UPDATE alerts_triggered_deliveries d
		SET lease_until = ?, leased_by = ?, updated_at = now()
		FROM cte
		WHERE d.id = cte.id
		RETURNING d.*;
	`, limit, leaseUntil, leaseOwner).Scan(&leased).Error

	if err != nil {
		return nil, fmt.Errorf("alert_triggered_delivery lease batch: %w", err)
	}

	out := make([]*domain.AlertTriggeredDelivery, 0, len(leased))
	for i := range leased {
		out = append(out, leased[i].ToDomain())
	}
	return out, nil
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

func (r *Repository) MarkSentMany(ctx context.Context, ids []int, sentAt time.Time) error {
	if len(ids) == 0 {
		return nil
	}

	arr := pgtype.Int8Array{}
	if err := arr.Set(ids); err != nil {
		return fmt.Errorf("deliveries mark sent many: set array: %w", err)
	}

	// status='sent', sent_at=..., updated_at=...
	if err := r.DB.WithContext(ctx).
		Model(&model.AlertTriggeredDelivery{}).
		Where("id = ANY(?)", &arr).
		Updates(map[string]any{
			"status":     "sent",
			"sent_at":    sentAt,
			"updated_at": time.Now(),
		}).Error; err != nil {

		return fmt.Errorf("deliveries mark sent many: %w", err)
	}

	return nil
}

func (r *Repository) MarkFailedMany(ctx context.Context, ids []int, errText string) error {
	if len(ids) == 0 {
		return nil
	}

	arr := pgtype.Int8Array{}
	if err := arr.Set(ids); err != nil {
		return fmt.Errorf("deliveries mark failed many: set array: %w", err)
	}

	if err := r.DB.WithContext(ctx).
		Model(&model.AlertTriggeredDelivery{}).
		Where("id = ANY(?)", &arr).
		Updates(map[string]any{
			"status":     "failed",
			"last_error": errText,
			"updated_at": time.Now(),
		}).Error; err != nil {

		return fmt.Errorf("deliveries mark failed many: %w", err)
	}

	return nil
}
