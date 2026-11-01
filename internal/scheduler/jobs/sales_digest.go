package jobs

import (
	"context"
	"log/slog"

	"salehelp/internal/services/sales_digest"
	"salehelp/pkg/logger/sl"
)

// SalesDigestJob — суточный дайджест выкупов и возвратов.
// Запускается раз в сутки после SyncStatsDailyJob (например, в 07:00 МСК).
type SalesDigestJob struct {
	log *slog.Logger
	svc *sales_digest.Service
}

func NewSalesDigestJob(log *slog.Logger, svc *sales_digest.Service) *SalesDigestJob {
	return &SalesDigestJob{
		log: log.With(sl.Module("SalesDigestJob")),
		svc: svc,
	}
}

func (j *SalesDigestJob) Run(ctx context.Context) {
	const op = "job.sales_digest"
	if err := j.svc.Run(ctx); err != nil {
		j.log.Error("sales digest failed", "op", op, sl.Err(err))
		return
	}
	j.log.Info("sales digest done", "op", op)
}
