package pattern

import (
	"context"
	"fmt"
	"sync"
	"orion/internal/types"
	"time"
)

// Engine manages the lifecycle and usage of cognitive patterns
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

// Match searches for a pattern trigger that matches the current goal context
func (e *Engine) Match(ctx context.Context, goal string) (*Pattern, bool) {
	fmt.Printf("Pattern Engine: Searching for trigger match for: %s\n", goal)

	pattern, found := e.store.MatchTrigger(ctx, goal)
	if found {
		e.eb.Publish(types.Event{
			Type: "pattern.matched",
			Payload: map[string]string{"pattern_id": pattern.ID},
			CreatedAt: time.Now(),
		})
	}

	return pattern, found
}

// RecordSuccess updates the pattern metrics after a successful execution
func (e *Engine) RecordSuccess(ctx context.Context, patternID string) {
	e.confidence.Increase(patternID)
	e.store.IncrementUsage(ctx, patternID)

	e.eb.Publish(types.Event{
		Type: "pattern.confidence_changed",
		Payload: map[string]string{"pattern_id": patternID, "status": "increased"},
		CreatedAt: time.Now(),
	})
}
