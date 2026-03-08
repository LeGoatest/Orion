package agents

import (
	"context"
	"fmt"
	"orion/internal/types"
	"time"
)

// Supervisor manages the lifecycle of all registered agents
type Supervisor struct {
	registry *Registry
	eventBus *types.EventBus
}

func NewSupervisor(registry *Registry, eb *types.EventBus) *Supervisor {
	return &Supervisor{
		registry: registry,
		eventBus: eb,
	}
}

// StartAgents initializes all registered agents
func (s *Supervisor) StartAgents(ctx context.Context) error {
	fmt.Println("Agent Supervisor: Starting agents...")
	for _, agent := range s.registry.ListAgents() {
		s.eventBus.Publish(types.Event{
			Type:      "agent.started",
			Payload:   agent.Name(),
			CreatedAt: time.Now(),
		})
	}
	return nil
}

// StopAgents gracefully shuts down all agents
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
