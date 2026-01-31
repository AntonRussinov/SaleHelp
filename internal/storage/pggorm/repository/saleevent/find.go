package saleevent

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// ListUnnotifiedByShop возвращает события магазина, по которым ещё не отправлено уведомление (notified_at IS NULL).
func (r *Repository) ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.SaleEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	var rows []model.SaleEvent
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND notified_at IS NULL", shopID).
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sale_events list unnotified: %w", err)
	}
	out := make([]*domain.SaleEvent, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}

// GetMaxLastChangeDateByShop возвращает максимальную last_change_date по магазину (для пагинации WB API).
// Если записей нет, возвращает zero time — вызывающий код подставит начальную дату (например 2019-06-20).
func (r *Repository) GetMaxLastChangeDateByShop(ctx context.Context, shopID int) (time.Time, error) {
	var maxTime *time.Time
	err := r.DB.WithContext(ctx).
		Model(&model.SaleEvent{}).
		Select("MAX(last_change_date)").
		Where("shop_id = ?", shopID).
		Scan(&maxTime).Error
	if err != nil {
		return time.Time{}, fmt.Errorf("sale_events get max last_change_date: %w", err)
	}
	if maxTime == nil {
		return time.Time{}, nil
	}
	return *maxTime, nil
}
