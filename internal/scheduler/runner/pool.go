package runner

import (
	"context"
	"sync"
)

type Task func(ctx context.Context)

type Pool struct {
	concurrency int
}

func New(concurrency int) *Pool {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &Pool{concurrency: concurrency}
}

func (p *Pool) RunAll(ctx context.Context, tasks []Task) {
	sem := make(chan struct{}, p.concurrency)
	var wg sync.WaitGroup

	for _, t := range tasks {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		default:
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(task Task) {
			defer wg.Done()
			defer func() { <-sem }()
			task(ctx)
		}(t)
	}

	wg.Wait()
}
