package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Scheduler manages job dispatching to the worker pool
type Scheduler struct {
	pool *WorkerPool
}

// NewScheduler creates a new Scheduler
func NewScheduler(pool *WorkerPool) *Scheduler {
	return &Scheduler{
		pool: pool,
	}
}

// Schedule submits a job to be executed by the worker pool
func (s *Scheduler) Schedule(ctx context.Context, jobType string, payload interface{}, handler func(context.Context, interface{}) error) (string, error) {
	jobID := uuid.New().String()
	job := Job{
		ID:        jobID,
		Type:      jobType,
		Payload:   payload,
		Handler:   handler,
		CreatedAt: time.Now(),
	}

	if err := s.pool.Submit(ctx, job); err != nil {
		return "", fmt.Errorf("failed to schedule job: %w", err)
	}

	return jobID, nil
}
