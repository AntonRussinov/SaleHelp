package shop_ai_reply_quota

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

// GetForMonth возвращает квоту для магазина за указанный месяц (первый день месяца).
// Если записи нет — nil, nil.
func (r *Repository) GetForMonth(ctx context.Context, shopID int, month time.Time) (*domain.ShopAIReplyQuota, error) {
	firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	var m model.ShopAIReplyQuota
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND month = ?", shopID, firstDay).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("shop_ai_reply_quotas get: %w", err)
	}
	return m.ToDomain(), nil
}
