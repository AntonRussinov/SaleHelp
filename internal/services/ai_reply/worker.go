package ai_reply

import (
	"context"
	"time"

	"salehelp/domain"
)

const (
	defaultBatchSize   = 10
	defaultItemTimeout = 60 * time.Second
)

// ProcessBatch обрабатывает пачку отзывов: для каждого создаёт context с таймаутом, вызывает ProcessOne, логирует ошибки без прерывания цикла.
func (s *Service) ProcessBatch(ctx context.Context, feedbacks []*domain.WBFeedback) {
	for _, f := range feedbacks {
		itemCtx, cancel := context.WithTimeout(ctx, defaultItemTimeout)
		err := s.ProcessOne(itemCtx, f)
		cancel()
		if err != nil {
			s.log.Error("ai_reply: process one failed", "feedback_id", f.ID, "shop_id", f.ShopID, "err", err)
		}
	}
}

// RunShop загружает до batchSize отзывов без ответа по магазину и обрабатывает их батчем.
func (s *Service) RunShop(ctx context.Context, shopID int, batchSize int) error {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}

	month := time.Now().UTC()
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	if err := s.quotaRepo.EnsureForMonth(ctx, shopID, month, s.includedLimit); err != nil {
		return err
	}
	quota, err := s.quotaRepo.GetForMonth(ctx, shopID, month)
	if err != nil {
		return err
	}
	if quota == nil || (quota.IncludedUsed >= quota.IncludedLimit && quota.PaidBalance <= 0) {
		return nil
	}

	feedbacks, err := s.feedbackRepo.ListUnrepliedByShop(ctx, shopID, batchSize)
	if err != nil {
		return err
	}
	if len(feedbacks) == 0 {
		return nil
	}
	s.log.Info("ai_reply: processing batch", "shop_id", shopID, "count", len(feedbacks))
	s.ProcessBatch(ctx, feedbacks)
	return nil
}
