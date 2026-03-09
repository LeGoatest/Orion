package agents

import (
	"context"
	"orion/internal/runtime"
	"time"
)

type Supervisor struct {
	Reg *Registry
	EB  *runtime.EventBus
}

func NewSupervisor(r *Registry, eb *runtime.EventBus) *Supervisor {
	return &Supervisor{
		Reg: r,
		EB:  eb,
	}
}

func (s *Supervisor) StartAgents(ctx context.Context) error {
	for _, a := range s.Reg.List() {
		s.EB.Publish(runtime.Event{
			Type:      "agent.started",
			Payload:   a.Name(),
			CreatedAt: time.Now(),
		})
	}
	return nil
}

func (s *Supervisor) StopAgents(ctx context.Context) error {
	return nil
}
