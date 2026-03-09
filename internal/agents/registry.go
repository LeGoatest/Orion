package agents

import "sync"

type Registry struct {
	mu sync.RWMutex
	m  map[string]AutonomousAgent
}

func NewRegistry() *Registry { return &Registry{m: make(map[string]AutonomousAgent)} }
func (r *Registry) RegisterAgent(a AutonomousAgent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[a.Name()] = a
}
func (r *Registry) GetByCap(c string) []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []AutonomousAgent
	for _, a := range r.m {
		for _, cap := range a.Capabilities() {
			if cap == c {
				res = append(res, a)
				break
			}
		}
	}
	return res
}
func (r *Registry) List() []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []AutonomousAgent
	for _, a := range r.m {
		res = append(res, a)
	}
	return res
}
