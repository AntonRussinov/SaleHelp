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

// DailyConsumptionByNmID возвращает средний дневной расход (выкупы - возвраты) по nm_id за период since..now.
// Ключ map — nm_id, значение — штук в день (может быть дробным).
func (r *Repository) DailyConsumptionByNmID(ctx context.Context, shopID int, since time.Time) (map[int64]float64, error) {
	type row struct {
		NmID     int64
		NetCount float64
	}
	var rows []row
	// buyout = +1, return = -1, partial_return = -1
	err := r.DB.WithContext(ctx).
		Raw(`
			SELECT nm_id,
				SUM(CASE
					WHEN event_type = 'buyout' THEN 1.0
					WHEN event_type IN ('return', 'partial_return') THEN -1.0
					ELSE 0
				END) / NULLIF(EXTRACT(EPOCH FROM (NOW() - ?)) / 86400, 0) AS net_count
			FROM sale_events
			WHERE shop_id = ? AND date >= ?
			GROUP BY nm_id
		`, since, shopID, since).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sale_events daily consumption by nm_id: %w", err)
	}
	out := make(map[int64]float64, len(rows))
	for _, rr := range rows {
		if rr.NetCount > 0 {
			out[rr.NmID] = rr.NetCount
		}
	}
	return out, nil
}

// GetMaxLastChangeDateByShop возвращает максимальную last_change_date по магазину (для пагинации WB API).
// Если записей нет, возвращает zero time — вызывающий код подставит начальную дату (например 2019-06-20).
func (r *Repository) GetMaxLastChangeDateByShop(ctx context.Context, shopID int) (time.Time, error) {
	type result struct {
		MaxDate *time.Time `gorm:"column:max_date"`
	}
	var res result
	err := r.DB.WithContext(ctx).
		Model(&model.SaleEvent{}).
		Select("MAX(last_change_date) as max_date").
		Where("shop_id = ?", shopID).
		Scan(&res).Error
	if err != nil {
		return time.Time{}, fmt.Errorf("sale_events get max last_change_date: %w", err)
	}
	if res.MaxDate == nil {
		return time.Time{}, nil
	}
	return *res.MaxDate, nil
}
