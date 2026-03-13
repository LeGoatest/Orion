package cognition

import (
	"context"
	"fmt"
	"time"
)

type GovernanceValidator struct{}

func (gv *GovernanceValidator) Validate(ctx context.Context, plan *ExecutionPlan, sm *SituationalModel) (*ValidatedExecutionPlan, error) {
	fmt.Printf("Governance: Validating plan for goal %s\n", plan.GoalID)

	// Simulated validation logic against SAGE-style categories
	disposition := "approved"
	reason := "All constraints satisfied."

	// Example security check: forbid dangerous tools in read-only mode
	for _, tool := range plan.Tools {
		for _, constraint := range sm.ConstraintSet {
			if constraint == "read-only" && tool == "filesystem_delete" {
				disposition = "rejected"
				reason = "Violation: Tool 'filesystem_delete' forbidden in 'read-only' mode."
				break
			}
		}
	}

	return &ValidatedExecutionPlan{
		Plan:             plan,
		Disposition:      disposition,
		ValidationReason: reason,
		ValidatedAt:      time.Now(),
	}, nil
}
