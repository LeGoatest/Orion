package agents

import (
	"context"
	"fmt"
	"orion/internal/runtime"
	"sync"
	"time"
)

type Supervisor struct {
	Reg    *Registry
	EB     *runtime.EventBus
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	status map[string]time.Time
}

func NewSupervisor(r *Registry, eb *runtime.EventBus) *Supervisor {
	return &Supervisor{
		Reg:    r,
		EB:     eb,
		status: make(map[string]time.Time),
	}
}

func (s *Supervisor) StartAgents(ctx context.Context) error {
	s.mu.Lock()
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.mu.Unlock()

	for _, a := range s.Reg.List() {
		s.startAgent(a)
	}

	go s.monitorLoop()
	return nil
}

func (s *Supervisor) startAgent(a AutonomousAgent) {
	fmt.Printf("Supervisor: Starting agent %s\n", a.Name())
	s.EB.Publish(runtime.Event{
		Type:      "agent.started",
		Payload:   a.Name(),
		CreatedAt: time.Now(),
	})

	s.mu.Lock()
	s.status[a.Name()] = time.Now()
	s.mu.Unlock()
}

func (s *Supervisor) monitorLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			for _, a := range s.Reg.List() {
				lastSeen, exists := s.status[a.Name()]
				if !exists || now.Sub(lastSeen) > 15*time.Second {
					fmt.Printf("Supervisor: Agent %s timed out or failed. Restarting...\n", a.Name())
					s.startAgent(a)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Supervisor) StopAgents(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func (s *Supervisor) Heartbeat(agentName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[agentName] = time.Now()
}
