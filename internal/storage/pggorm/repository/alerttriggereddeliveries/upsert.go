package alerttriggereddeliveries

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

// Upsert по UNIQUE(triggered_id, endpoint_id)
func (r *Repository) Upsert(ctx context.Context, d *domain.AlertTriggeredDelivery) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "triggered_id"},
				{Name: "endpoint_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"status",
				"sent_at",
				"attempts",
				"last_error",
				"lease_until",
				"leased_by",
				"updated_at",
			}),
		}).
		Create(m).Error; err != nil {

		return fmt.Errorf("alert_triggered_delivery upsert: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
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
