package wbfeedback

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// ListUnnotifiedByShop возвращает отзывы магазина, по которым ещё не отправлено уведомление.
func (r *Repository) ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.WBFeedback, error) {
	if limit <= 0 {
		limit = 100
	}
	var rows []model.WBFeedback
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND notified_at IS NULL", shopID).
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("wb_feedbacks list unnotified: %w", err)
	}
	out := make([]*domain.WBFeedback, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
