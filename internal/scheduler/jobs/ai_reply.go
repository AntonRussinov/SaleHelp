package jobs

import (
	"context"
	"log/slog"
)

// AISettingsRepo для списка магазинов с включённым AI.
type AISettingsRepo interface {
	ListActiveShopIDs(ctx context.Context, limit int) ([]int, error)
}

// AIReplyRunner сервис генерации AI-ответов.
type AIReplyRunner interface {
	RunShop(ctx context.Context, shopID int, batchSize int) error
}

// AIReplyJob — джоба: для каждого магазина с is_active=true запускает обработку батча отзывов.
type AIReplyJob struct {
	log    *slog.Logger
	settings AISettingsRepo
	runner   AIReplyRunner
	batchSize int
}

// NewAIReplyJob создаёт джобу.
func NewAIReplyJob(log *slog.Logger, settings AISettingsRepo, runner AIReplyRunner, batchSize int) *AIReplyJob {
	if batchSize <= 0 {
		batchSize = 10
	}
	return &AIReplyJob{log: log, settings: settings, runner: runner, batchSize: batchSize}
}

// Run выполняет один проход: список магазинов с AI → для каждого RunShop.
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
