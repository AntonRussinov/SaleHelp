package ai_reply

import (
	"context"
	"errors"
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

// RunPendingSends отправляет в WB сгенерированные ответы, у которых wb_sent_at IS NULL.
// Использует короткий контекст на каждую попытку: если rate-лимит заблокирован — пропускает и попробует на следующем запуске.
func (s *Service) RunPendingSends(ctx context.Context) error {
	const sendAttemptTimeout = 5 * time.Second
	const batchLimit = 10

	pending, err := s.replyRepo.ListPendingWBSend(ctx, batchLimit)
	if err != nil {
		return err
	}
	for _, p := range pending {
		shop, err := s.shopRepo.GetByID(ctx, p.ShopID)
		if err != nil || shop == nil || shop.APIToken == "" {
			s.log.Warn("ai_reply send: shop not found or no token", "shop_id", p.ShopID)
			continue
		}
		sendCtx, cancel := context.WithTimeout(ctx, sendAttemptTimeout)
		sendErr := s.wbFeedbackSend.AnswerFeedback(sendCtx, shop.APIToken, p.WBFeedbackID, p.AnswerText)
		cancel()
		if sendErr != nil {
			if errors.Is(sendErr, context.DeadlineExceeded) {
				s.log.Info("ai_reply send: rate limited, will retry next run", "reply_id", p.ReplyID)
			} else {
				s.log.Error("ai_reply send: failed", "reply_id", p.ReplyID, "wb_feedback_id", p.WBFeedbackID, "err", sendErr)
			}
			continue
		}
		if err := s.replyRepo.MarkWBSent(ctx, p.ReplyID, time.Now()); err != nil {
			s.log.Warn("ai_reply send: mark wb_sent_at failed", "reply_id", p.ReplyID, "err", err)
		}
	}
	return nil
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
