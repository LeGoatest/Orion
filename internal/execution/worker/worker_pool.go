package worker

import (
	"context"
	"fmt"
	"sync"
)

type Job interface {
	Execute(ctx context.Context) error
	ID() string
	Type() string
}

type WorkerPool struct {
	maxWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		jobQueue:   make(chan Job, 100),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.wg.Add(wp.maxWorkers)
	for i := 0; i < wp.maxWorkers; i++ {
		go wp.worker(ctx)
	}
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) worker(ctx context.Context) {
	defer wp.wg.Done()
	for {
		select {
		case job := <-wp.jobQueue:
			fmt.Printf("Worker executing job: %s (%s)\n", job.ID(), job.Type())
			if err := job.Execute(ctx); err != nil {
				fmt.Printf("Job %s failed: %v\n", job.ID(), err)
			}
		case <-ctx.Done():
			return
		}
	}
}
