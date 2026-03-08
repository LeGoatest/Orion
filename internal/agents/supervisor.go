package agents

import (
	"context"
	"fmt"
	"sync"
	"time"

	"orion/internal/types"
)

type Supervisor struct {
	mu           sync.RWMutex
	Registry     *Registry
	EventBus     *types.EventBus
	agentStates  map[string]*agentState
	stopChan     chan struct{}
}

type agentState struct {
	agent         AutonomousAgent
	lastHeartbeat time.Time
	status        string // "healthy", "failed", "restarting"
	cancel        context.CancelFunc
}

func NewSupervisor(r *Registry, eb *types.EventBus) *Supervisor {
	return &Supervisor{
		Registry:    r,
		EventBus:    eb,
		agentStates: make(map[string]*agentState),
		stopChan:    make(chan struct{}),
	}
}

func (s *Supervisor) StartAgents(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.Registry.ListAgents() {
		agentCtx, cancel := context.WithCancel(ctx)
		s.agentStates[a.Name()] = &agentState{
			agent:         a,
			lastHeartbeat: time.Now(),
			status:        "healthy",
			cancel:        cancel,
		}

		go s.runAgentHeartbeatLoop(agentCtx, a)
		s.EventBus.Publish(types.Event{Type: "agent.started", Payload: a.Name(), CreatedAt: time.Now()})
	}

	go s.monitorHealth(ctx)

	return nil
}

func (s *Supervisor) runAgentHeartbeatLoop(ctx context.Context, a AutonomousAgent) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.EventBus.Publish(types.Event{
				Type:      "agent.heartbeat",
				Payload:   a.Name(),
				CreatedAt: time.Now(),
			})
		case <-ctx.Done():
			return
		}
	}
}

func (s *Supervisor) monitorHealth(ctx context.Context) {
	// Subscribe to heartbeats
	ch := s.EventBus.Subscribe("agent.heartbeat")

	go func() {
		for {
			select {
			case event := <-ch:
				agentName, ok := event.Payload.(string)
				if ok {
					s.mu.Lock()
					if state, exists := s.agentStates[agentName]; exists {
						state.lastHeartbeat = time.Now()
						state.status = "healthy"
					}
					s.mu.Unlock()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Periodic check for dead agents
	checkTicker := time.NewTicker(10 * time.Second)
	defer checkTicker.Stop()

	for {
		select {
		case <-checkTicker.C:
			s.checkAgentHealth(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Supervisor) checkAgentHealth(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for name, state := range s.agentStates {
		if now.Sub(state.lastHeartbeat) > 15*time.Second && state.status == "healthy" {
			fmt.Printf("Supervisor: Agent %s failed health check (no heartbeat for 15s)\n", name)
			state.status = "failed"
			s.EventBus.Publish(types.Event{Type: "agent.failed", Payload: name, CreatedAt: now})

			// Restart Policy
			go s.restartAgent(ctx, name)
		}
	}
}

func (s *Supervisor) restartAgent(ctx context.Context, name string) {
	s.mu.Lock()
	state, exists := s.agentStates[name]
	if !exists {
		s.mu.Unlock()
		return
	}
	state.status = "restarting"
	state.cancel() // Cancel old heartbeat loop
	s.mu.Unlock()

	fmt.Printf("Supervisor: Restarting agent %s...\n", name)
	time.Sleep(1 * time.Second) // Small delay

	s.mu.Lock()
	newAgentCtx, cancel := context.WithCancel(ctx)
	state.cancel = cancel
	state.lastHeartbeat = time.Now()
	state.status = "healthy"
	go s.runAgentHeartbeatLoop(newAgentCtx, state.agent)
	s.mu.Unlock()

	s.EventBus.Publish(types.Event{Type: "agent.restarted", Payload: name, CreatedAt: time.Now()})
}

func (s *Supervisor) StopAgents(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, state := range s.agentStates {
		state.cancel()
		s.EventBus.Publish(types.Event{Type: "agent.stopped", Payload: name, CreatedAt: time.Now()})
	}
	return nil
}
