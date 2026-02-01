package di

import (
	"context"
	"fmt"
	"os"

	"salehelp/internal/scheduler"
	"salehelp/internal/scheduler/jobs"
)

func (prov *Provider) initScheduler(ctx context.Context) error {
	prov.scheduler = scheduler.New(prov.log)

	leasedBy := prov.cfg.InstanceID
	if leasedBy == "" {
		host, _ := os.Hostname()
		leasedBy = fmt.Sprintf("%s-%d", host, os.Getpid())
		prov.cfg.InstanceID = leasedBy
	}

	job := jobs.NewSyncActiveShopsJob(
		prov.log,
		jobs.SyncActiveShopsConfig{
			Spec:        prov.cfg.Scheduler.Spec,
			Concurrency: prov.cfg.Scheduler.Concurrency,
			ShopTimeout: prov.cfg.Scheduler.ShopTimeout,
			LeaseTTL:    prov.cfg.Scheduler.LeaseTTL,
			LeasedBy:    leasedBy,
			BatchSize:   prov.cfg.Scheduler.BatchSize,
		},
		prov.Repositories.Shop,
		prov.Services.Monitoring,
		prov.Repositories.ShopSyncLeases,
		prov.Services.Subscription,
	)

	notifyJob := jobs.NewNotifySubscriptionsJob(
		prov.log,
		jobs.NotifySubscriptionsConfig{
			Spec: prov.cfg.Scheduler.SubscriptionNotifySpec,
		},
		prov.Services.SubscriptionMonitor,
	)

	if err := prov.scheduler.Register(ctx, prov.cfg.Scheduler.SubscriptionNotifySpec, "notify_subscriptions", notifyJob.Run); err != nil {
		return err
	}

	stockReplenishJob := jobs.NewStockReplenishJob(
		prov.log,
		jobs.StockReplenishConfig{
			Spec: prov.cfg.Scheduler.StockReplenishSpec,
		},
		prov.Services.StockReplenish,
	)

	if prov.cfg.Scheduler.StockReplenishSpec != "" {
		if err := prov.scheduler.Register(ctx, prov.cfg.Scheduler.StockReplenishSpec, "stock_replenish", stockReplenishJob.Run); err != nil {
			return err
		}
	}

	// регистрируем задачу
	return prov.scheduler.Register(ctx, prov.cfg.Scheduler.Spec, "sync_active_shops", job.Run)
}

func (prov *Provider) StartScheduler(ctx context.Context) error {
	if prov.scheduler == nil {
		return nil
	}
	prov.scheduler.Start()
	return nil
}

func (prov *Provider) StopScheduler(ctx context.Context) error {
	if prov.scheduler == nil {
		return nil
	}
	return prov.scheduler.Stop(ctx)
}
