package cognition

import (
	"context"
	"orion/ent"
	"time"
	"github.com/google/uuid"
)

func Observe(ctx context.Context, client *ent.Client, intent string) (*NormalizedEvent, error) {
	goal, err := client.Goal.Create().
		SetID(uuid.New().String()).
		SetDescription(intent).
		SetStatus("pending").
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &NormalizedEvent{
		ID:        uuid.New().String(),
		GoalID:    goal.ID,
		Type:      "user_intent",
		Payload:   map[string]interface{}{"description": intent},
		Timestamp: time.Now(),
	}, nil
}
