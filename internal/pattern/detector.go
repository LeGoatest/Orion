package pattern

import (
	"context"
	"fmt"
	"orion/internal/types"
	"time"

	"github.com/google/uuid"
)

type Detector struct {
	eb    *types.EventBus
	store *Store
}

func NewDetector(eb *types.EventBus, store *Store) *Detector {
	return &Detector{
		eb:    eb,
		store: store,
	}
}

// Detect identifies recurring successful sequences from workspace history
func (d *Detector) Detect(ctx context.Context) {
	fmt.Println("Pattern Detector: Analyzing workspace history for successful sequences...")

	// Simulate detecting a successful pattern for "reindex repository"
	trigger := "reindex repository"
	if _, found := d.store.MatchTrigger(ctx, trigger); !found {
		p := &Pattern{
			ID:            uuid.New().String(),
			Trigger:       trigger,
			SolutionSteps: []string{"reindex_code", "update_callgraph", "refresh_embeddings"},
			Confidence:    0.5,
			UsageCount:    0,
			State:         StateActive,
			CreatedAt:     time.Now(),
		}

		if err := d.store.SavePattern(ctx, p); err == nil {
			d.eb.Publish(types.Event{
				Type:      "pattern.detected",
				Payload:   p,
				CreatedAt: time.Now(),
			})
		}
	}
}
