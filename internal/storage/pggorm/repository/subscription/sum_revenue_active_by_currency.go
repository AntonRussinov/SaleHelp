package subscription

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"time"
)

func (r *Repository) SumRevenueActiveByCurrency(ctx context.Context, now time.Time) ([]model.MoneyByCurrency, error) {
	var total int64
	err := r.DB.WithContext(ctx).Raw(
		"SELECT COALESCE(SUM(price_minor), 0) FROM subscriptions WHERE status IN ('active','trial') AND expires_at > ?", now,
	).Scan(&total).Error
	if err != nil {
		return nil, fmt.Errorf("sum revenue active: %w", err)
	}
	if total == 0 {
		return nil, nil
	}
	return []model.MoneyByCurrency{{CurrencyCode: "RUB", AmountMinor: total}}, nil
}

func (r *Repository) SumRevenueGraceByCurrency(ctx context.Context, now time.Time) ([]model.MoneyByCurrency, error) {
	graceCutoff := now.Add(-72 * time.Hour)
	var total int64
	err := r.DB.WithContext(ctx).Raw(
		"SELECT COALESCE(SUM(price_minor), 0) FROM subscriptions WHERE expires_at <= ? AND expires_at > ?", now, graceCutoff,
	).Scan(&total).Error
	if err != nil {
		return nil, fmt.Errorf("sum revenue grace: %w", err)
	}
	if total == 0 {
		return nil, nil
	}
	return []model.MoneyByCurrency{{CurrencyCode: "RUB", AmountMinor: total}}, nil
}
