package di

import (
	"context"
	"fmt"
	"os"

	"salehelp/internal/scheduler"
	"salehelp/internal/scheduler/jobs"
)

func (prov *Provider) initScheduler(ctx context.Context) error {
	// schedulerLoc выставляется в InitServices() (Europe/Moscow). Если по какой-то причине nil — fallback на Local внутри scheduler.New.
	prov.scheduler = scheduler.New(prov.log, prov.schedulerLoc)

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

	// AI Reply: генерация ответов на отзывы (каждые 10 мин или по своему расписанию)
	if prov.Services.AIReply != nil && prov.cfg.Scheduler.AIReplySpec != "" {
		aiReplyJob := jobs.NewAIReplyJob(prov.log, prov.Repositories.ShopAISettings, prov.Services.AIReply, 10)
		if err := prov.scheduler.Register(ctx, prov.cfg.Scheduler.AIReplySpec, "ai_reply", aiReplyJob.Run); err != nil {
			return err
		}
	}

	// ABC-анализ: еженедельный дайджест (воскресенье 23:00)
	abcWeeklySpec := prov.cfg.ABCAnalysis.WeeklySchedule
	if abcWeeklySpec == "" {
		abcWeeklySpec = "0 23 * * 0"
	}
	if prov.Services.ABCAnalysis != nil && abcWeeklySpec != "" {
		minSales := prov.cfg.ABCAnalysis.MinSalesThreshold
		if minSales <= 0 {
			minSales = 10
		}
		abcWeeklyJob := jobs.NewABCAnalysisJob(
			prov.log,
			jobs.ABCJobConfig{PeriodDays: 7, MinSalesThreshold: minSales, IsWeekly: true},
			prov.Repositories.Shop,
			prov.Services.ABCAnalysis,
			prov.Repositories.ShopABCSettings,
			prov.Services.Subscription,
			prov.Repositories.NotificationEndpoint,
			prov.TelegramAdapter(),
			prov.Repositories.Product,
		)
		if err := prov.scheduler.Register(ctx, abcWeeklySpec, "abc_weekly", abcWeeklyJob.Run); err != nil {
			return err
		}
	}

	// ABC-анализ: месячный отчёт (1-го числа в 01:00)
	abcMonthlySpec := prov.cfg.ABCAnalysis.MonthlySchedule
	if abcMonthlySpec == "" {
		abcMonthlySpec = "0 1 1 * *"
	}
	if prov.Services.ABCAnalysis != nil && abcMonthlySpec != "" {
		minSales := prov.cfg.ABCAnalysis.MinSalesThreshold
		if minSales <= 0 {
			minSales = 10
		}
		abcMonthlyJob := jobs.NewABCAnalysisJob(
			prov.log,
			jobs.ABCJobConfig{PeriodDays: 30, MinSalesThreshold: minSales, IsWeekly: false},
			prov.Repositories.Shop,
			prov.Services.ABCAnalysis,
			prov.Repositories.ShopABCSettings,
			prov.Services.Subscription,
			prov.Repositories.NotificationEndpoint,
			prov.TelegramAdapter(),
			prov.Repositories.Product,
		)
		if err := prov.scheduler.Register(ctx, abcMonthlySpec, "abc_monthly", abcMonthlyJob.Run); err != nil {
			return err
		}
	}

	// Дайджест отказов (не забрали в ПВЗ): ежедневно по расписанию из конфига (по умолчанию 07:00 МСК).
	if prov.Services.RefusalDigest != nil {
		refusalDigestSpec := prov.cfg.Scheduler.RefusalDigestSpec
		if refusalDigestSpec == "" {
			refusalDigestSpec = "0 7 * * *"
		}
		refusalDigestJob := jobs.NewRefusalDigestJob(prov.log, prov.Services.RefusalDigest)
		if err := prov.scheduler.Register(ctx, refusalDigestSpec, "refusal_digest", refusalDigestJob.Run); err != nil {
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
