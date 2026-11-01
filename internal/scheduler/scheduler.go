package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type jobMeta struct {
	name string
	spec string
}

type Scheduler struct {
	c    *cron.Cron
	log  Logger
	loc  *time.Location
	jobs map[cron.EntryID]jobMeta
}

func New(log Logger, loc *time.Location) *Scheduler {
	if loc == nil {
		loc = time.Local
	}
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	return &Scheduler{
		c:    c,
		log:  log,
		loc:  loc,
		jobs: make(map[cron.EntryID]jobMeta),
	}
}

func (s *Scheduler) Register(ctx context.Context, spec, name string, fn func(ctx context.Context)) error {
	id, err := s.c.AddFunc(spec, func() { fn(ctx) })
	if err != nil {
		return fmt.Errorf("register job %s (%s): %w", name, spec, err)
	}
	s.jobs[id] = jobMeta{name: name, spec: spec}
	s.log.Info("scheduler: registered job", "name", name, "spec", spec)
	return nil
}

func (s *Scheduler) Start() {
	s.c.Start()
	s.log.Info("scheduler: started", "location", s.loc.String())
	for id, meta := range s.jobs {
		e := s.c.Entry(id)
		// Next = zero, если расписание не может быть рассчитано (или уже остановлено).
		s.log.Info("scheduler: next run", "name", meta.name, "spec", meta.spec, "next", e.Next.Format(time.RFC3339))
	}
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
