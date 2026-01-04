package jobs

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"salehelp/internal/scheduler/runner"
)

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type ShopRepo interface {
	ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error)
}

type MonitoringService interface {
	SyncShop(ctx context.Context, shopID int) error
}

type ShopLeaseLocker interface {
	TryLease(ctx context.Context, shopID int, leasedBy string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, shopID int, leasedBy string) error
}

type SyncActiveShopsConfig struct {
	Spec        string
	Concurrency int
	ShopTimeout time.Duration
	LeaseTTL    time.Duration
	LeasedBy    string
	BatchSize   int
}

type SyncActiveShopsJob struct {
	log   Logger
	cfg   SyncActiveShopsConfig
	shops ShopRepo
	mon   MonitoringService
	lock  ShopLeaseLocker
	pool  *runner.Pool

	running atomic.Bool
}

var (
	okCount       int64
	failCount     int64
	skipCount     int64
	leaseErrCount int64
)

func NewSyncActiveShopsJob(
	log Logger,
	cfg SyncActiveShopsConfig,
	shops ShopRepo,
	mon MonitoringService,
	lock ShopLeaseLocker,
) *SyncActiveShopsJob {
	if cfg.LeasedBy == "" {
		host, _ := os.Hostname()
		cfg.LeasedBy = fmt.Sprintf("%s-%d", host, os.Getpid())
	}
	return &SyncActiveShopsJob{
		log:   log,
		cfg:   cfg,
		shops: shops,
		mon:   mon,
		lock:  lock,
		pool:  runner.New(cfg.Concurrency),
	}
}

func (j *SyncActiveShopsJob) Run(ctx context.Context) {
	if !j.running.CompareAndSwap(false, true) {
		j.log.Warn("sync_active_shops: previous run still active, skipping")
		return
	}
	defer j.running.Store(false)

	start := time.Now()
	j.log.Info("sync_active_shops: started")

	offset := 0
	total := 0

	for {
		select {
		case <-ctx.Done():
			j.log.Warn("sync_active_shops: canceled", "err", ctx.Err())
			j.log.Info("sync_active_shops: finished",
				"duration", time.Since(start).String(),
				"shops", total,
			)
			return
		default:
		}

		ids, err := j.shops.ListActiveIDs(ctx, j.cfg.BatchSize, offset)
		if err != nil {
			j.log.Error("sync_active_shops: list active shops failed", "err", err)
			break
		}
		if len(ids) == 0 {
			break
		}

		total += len(ids)

		tasks := make([]runner.Task, 0, len(ids))
		for _, shopID := range ids {
			id := shopID

			tasks = append(tasks, func(parent context.Context) {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				leased, err := j.lock.TryLease(parent, id, j.cfg.LeasedBy, j.cfg.LeaseTTL)
				if err != nil {
					j.log.Error("sync_shop: try lease failed", "shop_id", id, "err", err)
					return
				}
				if !leased {
					atomic.AddInt64(&skipCount, 1)
					return
				}
				defer func() { _ = j.lock.Release(parent, id, j.cfg.LeasedBy) }()

				shopCtx, cancel := context.WithTimeout(parent, j.cfg.ShopTimeout)
				defer cancel()

				if err := j.mon.SyncShop(shopCtx, id); err != nil {
					atomic.AddInt64(&failCount, 1)
					j.log.Error("sync_shop: failed", "shop_id", id, "err", err)
					return
				}

				atomic.AddInt64(&okCount, 1)
				j.log.Info("sync_shop: ok", "shop_id", id)
			})
		}

		j.pool.RunAll(ctx, tasks)

		offset += j.cfg.BatchSize
	}

	j.log.Info("sync_active_shops: finished",
		"duration", time.Since(start).String(),
		"shops", total,
	)

	j.log.Info("sync_active_shops: stats",
		"ok", okCount,
		"failed", failCount,
		"skipped", skipCount,
		"lease_errors", leaseErrCount,
	)

}
