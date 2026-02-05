package orderrefusal

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// ListUnnotifiedByShop возвращает отказы магазина, по которым ещё не отправлено уведомление (notified_at IS NULL).
func (r *Repository) ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.OrderRefusal, error) {
	if limit <= 0 {
		limit = 100
	}
	var rows []model.OrderRefusal
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND notified_at IS NULL", shopID).
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("order_refusals list unnotified: %w", err)
	}
	out := make([]*domain.OrderRefusal, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
