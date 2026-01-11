package jobs

import (
	"context"
	"log/slog"
	"salehelp/pkg/logger/sl"
)

type SubscriptionMonitor interface {
	Run(ctx context.Context) error
}

type NotifySubscriptionsConfig struct {
	Spec string
}

type NotifySubscriptionsJob struct {
	log     *slog.Logger
	cfg     NotifySubscriptionsConfig
	monitor SubscriptionMonitor
}

func NewNotifySubscriptionsJob(
	log *slog.Logger,
	cfg NotifySubscriptionsConfig,
	monitor SubscriptionMonitor,
) *NotifySubscriptionsJob {
	return &NotifySubscriptionsJob{
		log:     log.With(sl.Module("NotifySubscriptionsJob")),
		cfg:     cfg,
		monitor: monitor,
	}
}

func (j *NotifySubscriptionsJob) Run(ctx context.Context) {
	const op = "job.notify_subscriptions"

	if err := j.monitor.Run(ctx); err != nil {
		j.log.Error("notify subscriptions failed", op, sl.Err(err))
		return
	}

	j.log.Info("notify subscriptions done", "op", op)
}
