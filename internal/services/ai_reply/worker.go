package ai_reply

import (
	"context"
	"time"

	"salehelp/domain"
)

const (
	defaultBatchSize   = 10
	defaultItemTimeout = 25 * time.Second
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
