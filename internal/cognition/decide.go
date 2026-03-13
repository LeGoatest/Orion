package cognition

import (
	"context"
	"fmt"
)

func Decide(ctx context.Context, sm *SituationalModel) (*ExecutionPlan, error) {
	fmt.Printf("Cognition: Deciding on plan for goal %s\n", sm.GoalID)
	fmt.Printf("Orientation Summary: %s\n", sm.OrientationSummary)

	// Deterministic rule-based planner using the SituationalModel
	plan := &ExecutionPlan{
		GoalID: sm.GoalID,
		Steps:  []string{fmt.Sprintf("Rule-based execution for: %s", sm.OrientationSummary)},
		Tools:  sm.CapabilityCandidates,
	}

	return plan, nil
}
