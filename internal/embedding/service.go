package embedding

import (
	"context"
	"fmt"
	"sync"
	"orion/internal/types"
)

// Service coordinates embedding generation and persistence
type Service struct {
	provider Provider
	eb       *types.EventBus
	queue    chan Job
	wg       sync.WaitGroup
}

// NewService creates a new embedding service
func NewService(provider Provider, eb *types.EventBus) *Service {
	return &Service{
		provider: provider,
		eb:       eb,
		queue:    make(chan Job, 100),
	}
}

// Start initiates the background embedding worker
func (s *Service) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case job := <-s.queue:
				s.processJob(ctx, job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Service) processJob(ctx context.Context, job Job) {
	s.eb.Publish(types.Event{
		Type: "embedding.requested",
		Payload: map[string]string{"job_id": job.ID()},
	})

	err := job.Execute(ctx, s.provider)
	if err != nil {
		s.eb.Publish(types.Event{
			Type: "embedding.failed",
			Payload: map[string]string{"job_id": job.ID(), "error": err.Error()},
		})
		return
	}

	s.eb.Publish(types.Event{
		Type: "embedding.completed",
		Payload: map[string]string{"job_id": job.ID()},
	})
}

// Submit adds a job to the queue
func (s *Service) Submit(job Job) {
	s.queue <- job
}
