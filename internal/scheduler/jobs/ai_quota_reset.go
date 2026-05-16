package jobs

import (
	"context"
	"log/slog"
	"time"
)

type SubscriptionQuotaResetter interface {
	ResetMonthlyQuota(ctx context.Context, now time.Time) error
}

type AIQuotaResetJob struct {
	log   *slog.Logger
	subs  SubscriptionQuotaResetter
}

func NewAIQuotaResetJob(log *slog.Logger, subs SubscriptionQuotaResetter) *AIQuotaResetJob {
	return &AIQuotaResetJob{log: log, subs: subs}
}

func (j *AIQuotaResetJob) Run(ctx context.Context) {
	if err := j.subs.ResetMonthlyQuota(ctx, time.Now()); err != nil {
		j.log.Error("ai_quota_reset: failed", "err", err)
		return
	}
	j.log.Info("ai_quota_reset: monthly AI quota reset complete")
}
