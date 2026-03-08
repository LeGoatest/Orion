package agents

import (
	"context"
	"fmt"
	"orion/internal/types"
	"sync"
	"time"
)

type AgentState struct {
	Name      string
	Status    string
	Restarts  int
	LastActive time.Time
}

// Supervisor manages the lifecycle and tracks state of all registered agents
type Supervisor struct {
	registry *Registry
	eventBus *types.EventBus
	mu       sync.RWMutex
	states   map[string]*AgentState
}

func NewSupervisor(registry *Registry, eb *types.EventBus) *Supervisor {
	return &Supervisor{
		registry: registry,
		eventBus: eb,
		states:   make(map[string]*AgentState),
	}
}

func (s *Supervisor) StartAgents(ctx context.Context) error {
	fmt.Println("Agent Supervisor: Starting and monitoring agents...")
	for _, agent := range s.registry.ListAgents() {
		s.mu.Lock()
		s.states[agent.Name()] = &AgentState{Name: agent.Name(), Status: "started", LastActive: time.Now()}
		s.mu.Unlock()

		s.eventBus.Publish(types.Event{
			Type:      "agent.started",
			Payload:   agent.Name(),
			CreatedAt: time.Now(),
		})
	}
	return nil
}

func (s *Supervisor) StopAgents(ctx context.Context) error {
	fmt.Println("Agent Supervisor: Stopping agents...")
	for _, agent := range s.registry.ListAgents() {
		s.eventBus.Publish(types.Event{
			Type:      "agent.stopped",
			Payload:   agent.Name(),
			CreatedAt: time.Now(),
		})
	}
	return nil
}

func (s *Supervisor) RecordFailure(agentName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, ok := s.states[agentName]; ok {
		state.Status = "failed"
		state.Restarts++
		s.eventBus.Publish(types.Event{Type: "agent.failed", Payload: agentName, CreatedAt: time.Now()})
	}
}
