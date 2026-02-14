package shop_ai_reply_quota

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/ai"

	"gorm.io/gorm"
)

// EnsureForMonth создаёт запись квоты на месяц, если её ещё нет. includedLimit берётся из плана подписки.
func (r *Repository) EnsureForMonth(ctx context.Context, shopID int, month time.Time, includedLimit int) error {
	firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	var m model.ShopAIReplyQuota
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND month = ?", shopID, firstDay).
		First(&m).Error
	if err == nil {
		return nil // уже есть
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("shop_ai_reply_quotas ensure: %w", err)
	}
	m = model.ShopAIReplyQuota{
		ShopID:        shopID,
		Month:         firstDay,
		IncludedLimit: includedLimit,
		IncludedUsed:  0,
		PaidBalance:   0,
		PaidUsed:      0,
	}
	return r.DB.WithContext(ctx).Create(&m).Error
}

// DebitOne списывает один ответ: сначала из included, иначе из paid_balance. Вызывать в транзакции вместе с записью в ledger.
// Возвращает source ("included" или "paid") или ошибку при отсутствии квоты.
func (r *Repository) DebitOne(ctx context.Context, shopID int, month time.Time, feedbackID int) (source string, err error) {
	firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	var m model.ShopAIReplyQuota
	err = r.DB.WithContext(ctx).
		Where("shop_id = ? AND month = ?", shopID, firstDay).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", ai.ErrNoQuota
		}
		return "", fmt.Errorf("shop_ai_reply_quotas get: %w", err)
	}
	if m.IncludedUsed < m.IncludedLimit {
		m.IncludedUsed++
		now := time.Now()
		if err := r.DB.WithContext(ctx).Model(&m).Updates(map[string]any{
			"included_used": m.IncludedUsed,
			"updated_at":    now,
		}).Error; err != nil {
			return "", fmt.Errorf("shop_ai_reply_quotas update included: %w", err)
		}
		return domain.LedgerSourceIncluded, nil
	}
	if m.PaidBalance > 0 {
		m.PaidBalance--
		m.PaidUsed++
		if err := r.DB.WithContext(ctx).Model(&m).Updates(map[string]any{
			"paid_balance": m.PaidBalance,
			"paid_used":    m.PaidUsed,
			"updated_at":   time.Now(),
		}).Error; err != nil {
			return "", fmt.Errorf("shop_ai_reply_quotas update paid: %w", err)
		}
		return domain.LedgerSourcePaid, nil
	}
	return "", ai.ErrNoQuota
}
