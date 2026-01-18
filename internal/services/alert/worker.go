package alert

import (
	"context"
	"log/slog"
	"time"
)

type ShopsRepo interface {
	ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error)
}

type Worker struct {
	shopsRepo ShopsRepo
	engine    *AlertEngine
	log       *slog.Logger

	interval  time.Duration
	batchSize int
}

func NewWorker(log *slog.Logger, shopsRepo ShopsRepo, engine *AlertEngine, interval time.Duration) *Worker {
	if interval <= 0 {
		interval = 3 * time.Minute
	}
	return &Worker{
		shopsRepo: shopsRepo,
		engine:    engine,
		log:       log,
		interval:  interval,
		batchSize: 200,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.log.Info("ALERT_WORKER started")
	defer w.log.Info("ALERT_WORKER stopped")

	if w == nil || w.engine == nil {
		w.log.Error("alert worker: engine is nil")
		return
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// можно сделать первый прогон сразу
	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *Worker) runOnce(ctx context.Context) {
	if w.engine == nil {
		w.log.Error("alert engine is nil")
		return
	}
	offset := 0

	for {
		ids, err := w.shopsRepo.ListActiveIDs(ctx, w.batchSize, offset)
		if err != nil {
			w.log.Error("failed to list active shop IDs", "error", err)
			return
		}

		if len(ids) == 0 {
			return
		}

		for _, shopID := range ids {
			err = w.engine.RunShop(ctx, shopID) // лучше логировать err, но не стопать весь батч
			if err != nil {
				w.log.Error("failed to run alert engine for shop", "shopID", shopID, "error", err)
			}
		}

		offset += len(ids)
	}
}
