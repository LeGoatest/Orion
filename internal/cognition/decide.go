package cognition

import (
	"context"
	"fmt"
)

func Decide(ctx context.Context, sm *SituationalModel) (*ExecutionPlan, error) {
	// Deterministic rule-based planner
	plan := &ExecutionPlan{
		GoalID: sm.Goal.ID,
		Steps:  []string{fmt.Sprintf("Execute logic for: %s", sm.Goal.Description)},
		Tools:  []string{"logger"},
	}
	return plan, nil
}
