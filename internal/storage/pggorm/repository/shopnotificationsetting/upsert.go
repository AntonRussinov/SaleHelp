package shopnotificationsetting

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) Upsert(ctx context.Context, shopID int, eventKey string, enabled bool, now time.Time) error {
	if shopID <= 0 {
		return fmt.Errorf("shop_notification_setting upsert: invalid shop_id=%d", shopID)
	}
	if eventKey == "" {
		return fmt.Errorf("shop_notification_setting upsert: empty event_key")
	}
	if now.IsZero() {
		now = time.Now()
	}

	// NOTE: keep SQL explicit to control updated_at.
	q := `
INSERT INTO shop_notification_settings (shop_id, event_key, enabled, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (shop_id, event_key)
DO UPDATE SET enabled = EXCLUDED.enabled, updated_at = EXCLUDED.updated_at
`
	if err := r.DB.WithContext(ctx).Exec(q, shopID, eventKey, enabled, now, now).Error; err != nil {
		return fmt.Errorf("shop_notification_setting upsert: %w", err)
	}
	return nil
}
