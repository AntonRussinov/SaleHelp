// internal/storage/pggorm/repository/subscription/sum_revenue_active_by_currency.go
package subscription

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"time"
)

func (r *Repository) SumRevenueActiveByCurrency(ctx context.Context, now time.Time) ([]model.MoneyByCurrency, error) {
	var rows []model.MoneyByCurrency

	// Берём по каждому shop одну "текущую активную" подписку: с максимальным valid_to
	// Потом суммируем base + addons по валюте подписки.
	q := `
WITH active_sub AS (
  SELECT DISTINCT ON (s.shop_id)
    s.id, s.shop_id, s.currency_code, s.base_amount_minor
  FROM subscriptions s
  WHERE s.status = 'active' AND s.valid_to > $1
  ORDER BY s.shop_id, s.valid_to DESC, s.id DESC
),
addon_sum AS (
  SELECT a.subscription_id, COALESCE(SUM(a.amount_minor), 0) AS addons_minor
  FROM subscription_addons a
  GROUP BY a.subscription_id
)
SELECT
  a.currency_code,
  SUM(a.base_amount_minor + COALESCE(x.addons_minor, 0)) AS amount_minor,
  COUNT(*) AS cnt
FROM active_sub a
LEFT JOIN addon_sum x ON x.subscription_id = a.id
GROUP BY a.currency_code
ORDER BY a.currency_code;
`

	if err := r.DB.WithContext(ctx).Raw(q, now).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("sum revenue active by currency: %w", err)
	}

	return rows, nil
}

func (r *Repository) SumRevenueGraceByCurrency(ctx context.Context, now time.Time) ([]model.MoneyByCurrency, error) {
	var rows []model.MoneyByCurrency

	q := `
WITH grace_sub AS (
  SELECT DISTINCT ON (s.shop_id)
    s.id, s.shop_id, s.currency_code, s.base_amount_minor
  FROM subscriptions s
  WHERE s.status = 'active'
    AND s.valid_to <= $1
    AND (s.valid_to + interval '3 days') > $1
  ORDER BY s.shop_id, s.valid_to DESC, s.id DESC
),
addon_sum AS (
  SELECT a.subscription_id, COALESCE(SUM(a.amount_minor), 0) AS addons_minor
  FROM subscription_addons a
  GROUP BY a.subscription_id
)
SELECT
  g.currency_code,
  SUM(g.base_amount_minor + COALESCE(x.addons_minor, 0)) AS amount_minor,
  COUNT(*) AS cnt
FROM grace_sub g
LEFT JOIN addon_sum x ON x.subscription_id = g.id
GROUP BY g.currency_code
ORDER BY g.currency_code;
`

	if err := r.DB.WithContext(ctx).Raw(q, now).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("sum revenue grace by currency: %w", err)
	}

	return rows, nil
}
