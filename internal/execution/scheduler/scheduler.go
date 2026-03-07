package scheduler

import (
	"context"
	"fmt"
	"orion/internal/execution/worker"
	"orion/internal/types"
)

type Scheduler struct {
	eventBus *types.EventBus
	pool     *worker.WorkerPool
}

func NewScheduler(eb *types.EventBus, wp *worker.WorkerPool) *Scheduler {
	return &Scheduler{
		eventBus: eb,
		pool:     wp,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	fmt.Println("Scheduler started.")
	// Periodic job polling logic would go here
}

func (s *Scheduler) Schedule(job worker.Job) {
	s.pool.Submit(job)
}
