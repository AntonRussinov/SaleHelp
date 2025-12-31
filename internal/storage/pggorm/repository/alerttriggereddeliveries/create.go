package alerttriggereddeliveries

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, d *domain.AlertTriggeredDelivery) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("alert_triggered_delivery create: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) CreateDeliveriesForTriggered(ctx context.Context, triggeredID int, shopID int) error {
	q := `
		INSERT INTO alerts_triggered_deliveries (triggered_id, endpoint_id)
		SELECT ?, e.id
		FROM shop_notification_endpoints e
		WHERE e.shop_id = ?
		  AND e.enabled = TRUE
		ON CONFLICT (triggered_id, endpoint_id) DO NOTHING;
	`

	if err := r.DB.WithContext(ctx).Exec(q, triggeredID, shopID).Error; err != nil {
		return fmt.Errorf("fanout create deliveries (triggered_id=%d shop_id=%d): %w", triggeredID, shopID, err)
	}
	return nil
}
