package jobs

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"salehelp/pkg/logger/sl"
)

// DailyStatsMonitoringService — интерфейс для ежедневного синка продаж и отказов.
type DailyStatsMonitoringService interface {
	SyncShopDailyStats(ctx context.Context, shopID int) error
}

// SyncStatsDailyJob — ежедневный синк продаж (supplier/sales) и отказов (supplier/orders) из WB.
// Вынесен из fast sync для снижения расхода платных WB API вызовов.
// Запускается раз в сутки до дайджестов (например, в 06:00 МСК).
type SyncStatsDailyJob struct {
	log         *slog.Logger
	shops       ShopRepo
	mon         DailyStatsMonitoringService
	concurrency int
	shopTimeout time.Duration

	running atomic.Bool
}

func NewSyncStatsDailyJob(
	log *slog.Logger,
	shops ShopRepo,
	mon DailyStatsMonitoringService,
	concurrency int,
	shopTimeout time.Duration,
) *SyncStatsDailyJob {
	if concurrency <= 0 {
		concurrency = 3
	}
	if shopTimeout <= 0 {
		shopTimeout = 5 * time.Minute
	}
	return &SyncStatsDailyJob{
		log:         log.With(sl.Module("SyncStatsDailyJob")),
		shops:       shops,
		mon:         mon,
		concurrency: concurrency,
		shopTimeout: shopTimeout,
	}
}

func (j *SyncStatsDailyJob) Run(ctx context.Context) {
	if !j.running.CompareAndSwap(false, true) {
		j.log.Warn("sync_stats_daily: previous run still active, skipping")
		return
	}
	defer j.running.Store(false)

	start := time.Now()
	j.log.Info("sync_stats_daily: started")

	const batchSize = 50
	offset := 0
	total := 0

	sem := make(chan struct{}, j.concurrency)

	for {
		select {
		case <-ctx.Done():
			j.log.Warn("sync_stats_daily: canceled", "err", ctx.Err())
			return
		default:
		}

		ids, err := j.shops.ListActiveIDs(ctx, batchSize, offset)
		if err != nil {
			j.log.Error("sync_stats_daily: list shops failed", "err", err)
			break
		}
		if len(ids) == 0 {
			break
		}

		total += len(ids)

		done := make(chan struct{}, len(ids))
		for _, shopID := range ids {
			id := shopID
			sem <- struct{}{}
			go func() {
				defer func() {
					<-sem
					done <- struct{}{}
				}()
				shopCtx, cancel := context.WithTimeout(ctx, j.shopTimeout)
				defer cancel()
				if err := j.mon.SyncShopDailyStats(shopCtx, id); err != nil {
					j.log.Error("sync_stats_daily: shop failed", "shop_id", id, sl.Err(err))
				}
			}()
		}
		for range ids {
			<-done
		}

		offset += batchSize
		if len(ids) < batchSize {
			break
		}
	}

	j.log.Info("sync_stats_daily: finished",
		"duration", time.Since(start).String(),
		"shops", total,
	)
}
