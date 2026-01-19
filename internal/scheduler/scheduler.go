package scheduler

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Scheduler struct {
	c   *cron.Cron
	log Logger
}

func New(log Logger) *Scheduler {
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	return &Scheduler{c: c, log: log}
}

func (s *Scheduler) Register(ctx context.Context, spec, name string, fn func(ctx context.Context)) error {
	_, err := s.c.AddFunc(spec, func() { fn(ctx) })
	if err != nil {
		return fmt.Errorf("register job %s (%s): %w", name, spec, err)
	}
	s.log.Info("scheduler: registered job", "name", name, "spec", spec)
	return nil
}

func (s *Scheduler) Start() {
	s.c.Start()
	s.log.Info("scheduler: started")
}

func (s *Scheduler) Stop(ctx context.Context) error {
	stopCtx := s.c.Stop()
	select {
	case <-stopCtx.Done():
		s.log.Info("scheduler: stopped")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
