package execution

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job represents a generic task to be executed
type Job struct {
	ID        string
	Type      string
	Payload   interface{}
	Handler   func(context.Context, interface{}) error
	CreatedAt time.Time
}

// WorkerPool manages background execution units
type WorkerPool struct {
	workerCount int
	jobQueue    chan Job
	wg          sync.WaitGroup
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan Job, queueSize),
	}
}

// Start spawns workers and begins job processing
func (wp *WorkerPool) Start(ctx context.Context) {
	wp.wg.Add(wp.workerCount)
	for i := 0; i < wp.workerCount; i++ {
		go wp.worker(ctx, i)
	}
}

// Submit adds a job to the queue
func (wp *WorkerPool) Submit(ctx context.Context, job Job) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("worker pool queue is full")
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()
	fmt.Printf("Worker %d starting\n", id)

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				fmt.Printf("Worker %d: queue closed, shutting down\n", id)
				return
			}
			fmt.Printf("Worker %d: executing job %s (%s)\n", id, job.ID, job.Type)
			if err := job.Handler(ctx, job.Payload); err != nil {
				fmt.Printf("Worker %d: job %s failed: %v\n", id, job.ID, err)
			}
		case <-ctx.Done():
			fmt.Printf("Worker %d: context cancelled, shutting down\n", id)
			return
		}
	}
}

// Wait blocks until all workers have stopped
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Stop closes the job queue and waits for workers to finish
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.Wait()
}
