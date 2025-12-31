package alert

import (
	"context"
	"time"
)

type ShopsRepo interface {
	ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error)
}

type Worker struct {
	shopsRepo ShopsRepo
	engine    *AlertEngine

	interval  time.Duration
	batchSize int
}

func NewWorker(shopsRepo ShopsRepo, engine *AlertEngine, interval time.Duration) *Worker {
	if interval <= 0 {
		interval = 3 * time.Minute
	}
	return &Worker{
		shopsRepo: shopsRepo,
		engine:    engine,
		interval:  interval,
		batchSize: 200,
	}
}

func (w *Worker) Run(ctx context.Context) {
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
	offset := 0

	for {
		ids, err := w.shopsRepo.ListActiveIDs(ctx, w.batchSize, offset)
		if err != nil {
			// тут лучше логировать; ошибку наружу не кидаем, чтобы цикл не умер
			return
		}
		if len(ids) == 0 {
			return
		}

		for _, shopID := range ids {
			_ = w.engine.RunShop(ctx, shopID) // лучше логировать err, но не стопать весь батч
		}

		offset += len(ids)
	}
}
