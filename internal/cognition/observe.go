package cognition

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Goal represents a user intent or sub-task
type Goal struct {
	ID          string
	Description string
	Status      string
	Priority    int
	CreatedAt   time.Time
}

// Observe captures user intent and creates a user_goal memory node
func (p *DefaultPipeline) Observe(ctx context.Context, input interface{}) (interface{}, error) {
	fmt.Println("Observing user intent...")

	intent, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("invalid input type for observe: expected string")
	}

	goal := &Goal{
		ID:          uuid.New().String(),
		Description: intent,
		Status:      "pending",
		Priority:    0,
		CreatedAt:   time.Now(),
	}

	// In a real implementation, this would be persisted to the database
	// and an 'observe' event would be emitted to the event bus.
	fmt.Printf("Goal created: %s (%s)\n", goal.Description, goal.ID)

	return goal, nil
}
