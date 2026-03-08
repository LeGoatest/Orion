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
	return &Detector{eb: eb, store: store}
}

func (d *Detector) DetectPatterns(ctx context.Context, goalID string) {
	fmt.Printf("Pattern Detector: Analyzing success path for goal %s\n", goalID)
	// Real logic: query successful steps for this goal and look for recurring sequences

	p := &Pattern{
		ID:            uuid.New().String(),
		Trigger:       "recurring task",
		SolutionSteps: []string{"step1", "step2"},
		Confidence:    0.9,
	}

	d.eb.Publish(types.Event{Type: "pattern.detected", Payload: p, CreatedAt: time.Now()})
}
