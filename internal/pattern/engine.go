package pattern

import (
	"context"
	"orion/internal/types"
	"sync"
	"time"
)

type State string

const (
	StateActive   State = "active"
	StateDegraded State = "degraded"
	StateArchived State = "archived"
)

type Pattern struct {
	ID            string    `json:"id"`
	Trigger       string    `json:"trigger"`
	SolutionSteps []string  `json:"solution_steps"`
	Confidence    float64   `json:"confidence"`
	UsageCount    int       `json:"usage_count"`
	State         State     `json:"state"`
	CreatedAt     time.Time `json:"created_at"`
}

type Engine struct {
	mu         sync.RWMutex
	store      *Store
	eb         *types.EventBus
	confidence *ConfidenceTracker
}

func NewEngine(store *Store, eb *types.EventBus) *Engine {
	return &Engine{
		store:      store,
		eb:         eb,
		confidence: NewConfidenceTracker(),
	}
}

func (e *Engine) Match(ctx context.Context, goal string) (*Pattern, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	pattern, found := e.store.MatchTrigger(ctx, goal)
	if found && pattern.Confidence > 0.8 {
		e.eb.Publish(types.Event{
			Type:      "pattern.matched",
			Payload:   pattern,
			CreatedAt: time.Now(),
		})
		return pattern, true
	}
	return nil, false
}

func (e *Engine) ExecutePattern(ctx context.Context, pattern *Pattern) error {
	e.eb.Publish(types.Event{
		Type:      "pattern.executed",
		Payload:   pattern.ID,
		CreatedAt: time.Now(),
	})
	// Actual execution logic would go here, probably calling back into cognition
	return nil
}
