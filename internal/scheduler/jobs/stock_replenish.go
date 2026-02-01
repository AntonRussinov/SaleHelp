package jobs

import (
	"context"
	"log/slog"

	"salehelp/internal/services/stock_replenish"
	"salehelp/pkg/logger/sl"
)

type StockReplenishConfig struct {
	Spec string
}

type StockReplenishJob struct {
	log   *slog.Logger
	cfg   StockReplenishConfig
	svc   *stock_replenish.Service
}

func NewStockReplenishJob(
	log *slog.Logger,
	cfg StockReplenishConfig,
	svc *stock_replenish.Service,
) *StockReplenishJob {
	return &StockReplenishJob{
		log: log.With(sl.Module("StockReplenishJob")),
		cfg: cfg,
		svc: svc,
	}
}

func (j *StockReplenishJob) Run(ctx context.Context) {
	const op = "job.stock_replenish"

	if err := j.svc.Run(ctx); err != nil {
		j.log.Error("stock replenish failed", "op", op, sl.Err(err))
		return
	}

	j.log.Info("stock replenish done", "op", op)
}
