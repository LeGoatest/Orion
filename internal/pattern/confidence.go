package pattern

import (
	"sync"
)

// ConfidenceTracker manages the confidence level of patterns
type ConfidenceTracker struct {
	mu          sync.RWMutex
	confidence map[string]float64
}

func NewConfidenceTracker() *ConfidenceTracker {
	return &ConfidenceTracker{
		confidence: make(map[string]float64),
	}
}

func (t *ConfidenceTracker) Increase(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.confidence[id] += 0.05
	if t.confidence[id] > 1.0 { t.confidence[id] = 1.0 }
}

func (t *ConfidenceTracker) Decrease(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.confidence[id] -= 0.1
	if t.confidence[id] < 0.0 { t.confidence[id] = 0.0 }
}
