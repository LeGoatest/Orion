package execution

import (
	"context"
	"fmt"
	"orion/internal/types"
	"sync"
)

type WorkerPool struct {
	max   int
	queue chan types.Job
	wg    sync.WaitGroup
}

func NewWorkerPool(m int) *WorkerPool {
	return &WorkerPool{
		max:   m,
		queue: make(chan types.Job, 100),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.wg.Add(wp.max)
	for i := 0; i < wp.max; i++ {
		go func(id int) {
			defer wp.wg.Done()
			for {
				select {
				case j := <-wp.queue:
					// Job is now a function in types.Job
					j(ctx)
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

func (wp *WorkerPool) Submit(j types.Job) {
	wp.queue <- j
}

func (wp *WorkerPool) Shutdown() {
	// Simple wait for goroutines to finish after context cancel
	wp.wg.Wait()
}

type Scheduler struct {
	pool *WorkerPool
}

func NewScheduler(p *WorkerPool) *Scheduler {
	return &Scheduler{pool: p}
}

func (s *Scheduler) Schedule(j types.Job) {
	s.pool.Submit(j)
}
