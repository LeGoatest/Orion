package cognition

import (
	"fmt"
	"sync"
)

// ExecutorRegistry manages the binding of algorithm IDs to concrete executors.
type ExecutorRegistry struct {
	mu         sync.RWMutex
	algorithms map[string]Algorithm
}

func NewExecutorRegistry() *ExecutorRegistry {
	r := &ExecutorRegistry{
		algorithms: make(map[string]Algorithm),
	}
	// Auto-register default algorithms
	r.Register(&ObserveAlgorithm{})
	r.Register(&OrientAlgorithm{})
	r.Register(&DecideAlgorithm{})
	return r
}

func (r *ExecutorRegistry) Register(algo Algorithm) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.algorithms[algo.ID()] = algo
}

func (r *ExecutorRegistry) Get(id string) (Algorithm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	algo, ok := r.algorithms[id]
	if !ok {
		return nil, fmt.Errorf("unknown algorithm ID: %s", id)
	}
	return algo, nil
}
