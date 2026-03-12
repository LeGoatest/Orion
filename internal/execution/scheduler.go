package execution

import (
	"context"
)

type Job func(ctx context.Context)

type Scheduler struct {
	pool *WorkerPool
}

func NewScheduler(pool *WorkerPool) *Scheduler {
	return &Scheduler{pool: pool}
}

func (s *Scheduler) Schedule(ctx context.Context, job Job) {
	s.pool.Submit(job)
}
