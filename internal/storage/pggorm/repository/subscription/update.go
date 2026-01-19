package subscription

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"time"
)

// ExpireActiveByShopID переводит все активные подписки магазина в expired.
// Возвращает кол-во обновлённых строк.
func (r *Repository) ExpireActiveByShopID(ctx context.Context, shopID int) (int64, error) {
	res := r.DB.WithContext(ctx).
		Model(&model.Subscription{}).
		Where("shop_id = ? AND status = ?", shopID, "active").
		Updates(map[string]any{
			"status":     "expired",
			"updated_at": time.Now(),
		})

	if res.Error != nil {
		return 0, fmt.Errorf("expire active subscriptions by shop_id: %w", res.Error)
	}

	return res.RowsAffected, nil
}
