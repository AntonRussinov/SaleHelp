package shop_ai_settings

import (
	"context"
	"fmt"
)

// ListActiveShopIDs возвращает ID магазинов с включёнными AI-ответами (is_active = true).
func (r *Repository) ListActiveShopIDs(ctx context.Context, limit int) ([]int, error) {
	if limit <= 0 {
		limit = 500
	}
	var ids []int
	err := r.DB.WithContext(ctx).
		Table("shop_ai_settings").
		Where("is_active = ?", true).
		Limit(limit).
		Pluck("shop_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("shop_ai_settings list active: %w", err)
	}
	return ids, nil
}
