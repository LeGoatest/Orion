package cognition

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Goal represents a minimal goal structure for the cognition package to avoid import cycle
type Goal struct {
	ID          string
	Description string
	Status      string
	CreatedAt   time.Time
}

// Observe implementation
func (p *DefaultPipeline) Observe(ctx context.Context, input interface{}) (interface{}, error) {
	fmt.Println("OODA-L: Phase: Observe")

	intent, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected string")
	}

	goal := &Goal{
		ID:          uuid.New().String(),
		Description: intent,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	fmt.Printf("User intent captured: %s (%s)\n", goal.Description, goal.ID)

	return goal, nil
}
