package shop_ai_reply_ledger

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// Create записывает операцию списания в аудит.
func (r *Repository) Create(ctx context.Context, d *domain.ShopAIReplyLedger) error {
	m := &model.ShopAIReplyLedger{
		ShopID:     d.ShopID,
		FeedbackID: d.FeedbackID,
		Month:      d.Month,
		Source:     d.Source,
		Delta:      d.Delta,
		Reason:     d.Reason,
	}
	err := r.DB.WithContext(ctx).Create(m).Error
	if err != nil {
		return fmt.Errorf("shop_ai_reply_ledger create: %w", err)
	}
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	return nil
}
