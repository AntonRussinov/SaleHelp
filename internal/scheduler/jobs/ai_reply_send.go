package jobs

import (
	"context"
	"log/slog"
)

// AIReplySendJob отправляет в WB сгенерированные ответы с wb_sent_at IS NULL.
// Запускается каждые 12 минут, что соответствует лимиту WB 5 запросов/час
// (POST /api/v1/feedbacks/answer). Один запуск = один ответ (burst=1 в rate limiter).
type AIReplySendJob struct {
	log    *slog.Logger
	sender AIReplySender
}

// NewAIReplySendJob создаёт джобу отправки.
func NewAIReplySendJob(log *slog.Logger, sender AIReplySender) *AIReplySendJob {
	return &AIReplySendJob{log: log, sender: sender}
}

// Run выполняет один проход: отправляет до batch=10 ответов, пропускает при rate limit.
func (j *AIReplySendJob) Run(ctx context.Context) {
	if err := j.sender.RunPendingSends(ctx); err != nil {
		j.log.Error("ai_reply_send job: failed", "err", err)
	}
}
