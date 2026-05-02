package jobs

import (
	"context"
	"log/slog"
)

// AISettingsRepo для списка магазинов с включённым AI.
type AISettingsRepo interface {
	ListActiveShopIDs(ctx context.Context, limit int) ([]int, error)
}

// AIReplyGenerator генерирует AI-ответы на отзывы (используется AIReplyJob).
type AIReplyGenerator interface {
	RunShop(ctx context.Context, shopID int, batchSize int) error
}

// AIReplySender отправляет сгенерированные ответы в WB (используется AIReplySendJob).
type AIReplySender interface {
	RunPendingSends(ctx context.Context) error
}

// AIReplyJob — генерирует AI-ответы для всех активных магазинов. Отправка в WB — задача AIReplySendJob.
type AIReplyJob struct {
	log       *slog.Logger
	settings  AISettingsRepo
	runner    AIReplyGenerator
	batchSize int
}

// NewAIReplyJob создаёт джобу генерации.
func NewAIReplyJob(log *slog.Logger, settings AISettingsRepo, runner AIReplyGenerator, batchSize int) *AIReplyJob {
	if batchSize <= 0 {
		batchSize = 10
	}
	return &AIReplyJob{log: log, settings: settings, runner: runner, batchSize: batchSize}
}

// Run генерирует ответы на отзывы для каждого активного магазина.
func (j *AIReplyJob) Run(ctx context.Context) {
	shopIDs, err := j.settings.ListActiveShopIDs(ctx, 200)
	if err != nil {
		j.log.Error("ai_reply job: list active shops", "err", err)
		return
	}
	for _, shopID := range shopIDs {
		if err := j.runner.RunShop(ctx, shopID, j.batchSize); err != nil {
			j.log.Error("ai_reply job: run shop failed", "shop_id", shopID, "err", err)
		}
	}
}
