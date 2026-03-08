package agents

import (
	"sync"
)

// Registry maintains a list of available agents and their metadata
type Registry struct {
	mu     sync.RWMutex
	agents map[string]AutonomousAgent
}

func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]AutonomousAgent),
	}
}

// RegisterAgent adds a new agent to the system
func (r *Registry) RegisterAgent(agent AutonomousAgent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[agent.Name()] = agent
}

// GetAgentsByCapability returns agents that support a specific capability
func (r *Registry) GetAgentsByCapability(capability string) []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []AutonomousAgent
	for _, a := range r.agents {
		for _, c := range a.Capabilities() {
			if c == capability {
				matches = append(matches, a)
				break
			}
		}
	}
	return matches
}

func (r *Registry) ListAgents() []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]AutonomousAgent, 0, len(r.agents))
	for _, a := range r.agents {
		list = append(list, a)
	}
	return list
}
