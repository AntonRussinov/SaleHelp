package jobs

import (
	"context"
	"log/slog"

	"salehelp/internal/services/refusal_digest"
	"salehelp/pkg/logger/sl"
)

// RefusalDigestJob — суточный дайджест отказов (заказы не забраны в ПВЗ).
// Запускается раз в сутки по расписанию (07:00 МСК).
type RefusalDigestJob struct {
	log *slog.Logger
	svc *refusal_digest.Service
}

func NewRefusalDigestJob(log *slog.Logger, svc *refusal_digest.Service) *RefusalDigestJob {
	return &RefusalDigestJob{
		log: log.With(sl.Module("RefusalDigestJob")),
		svc: svc,
	}
}

func (j *RefusalDigestJob) Run(ctx context.Context) {
	const op = "job.refusal_digest"
	if err := j.svc.Run(ctx); err != nil {
		j.log.Error("refusal digest failed", "op", op, sl.Err(err))
		return
	}
	j.log.Info("refusal digest done", "op", op)
}
