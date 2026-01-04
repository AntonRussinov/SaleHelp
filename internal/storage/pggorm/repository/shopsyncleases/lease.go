package shopsyncleases

import (
	"context"
	"time"
)

// TryLease пытается взять lease на shopID до now()+ttl.
// Возвращает true если lease получен.
func (r *Repository) TryLease(ctx context.Context, shopID int, leasedBy string, ttl time.Duration) (bool, error) {
	var ok int

	// Смысл:
	// - Если строки нет -> INSERT
	// - Если есть и lease истёк -> UPDATE
	// - Если есть и lease активен -> НИЧЕГО (ok останется 0)
	q := `
WITH upsert AS (
  INSERT INTO shop_sync_leases (shop_id, leased_by, lease_until)
  VALUES (?, ?, now() + (? * interval '1 second'))
  ON CONFLICT (shop_id) DO UPDATE
    SET leased_by   = EXCLUDED.leased_by,
        lease_until = EXCLUDED.lease_until
    WHERE shop_sync_leases.lease_until IS NULL
       OR shop_sync_leases.lease_until < now()
  RETURNING 1
)
SELECT 1 FROM upsert;
`
	seconds := int(ttl.Seconds())
	err := r.DB.WithContext(ctx).Raw(q, shopID, leasedBy, seconds).Scan(&ok).Error
	if err != nil {
		return false, err
	}
	return ok == 1, nil
}

func (r *Repository) Release(ctx context.Context, shopID int, leasedBy string) error {
	// Release только если мы владелец lease — чтобы случайно не снять чужой
	return r.DB.WithContext(ctx).
		Exec(`
			UPDATE shop_sync_leases
			SET leased_by = NULL, lease_until = NULL
			WHERE shop_id = ? AND leased_by = ?;
			`, shopID, leasedBy).Error
}
