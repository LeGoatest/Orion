package cognition

import (
	"context"
	"fmt"
)

func Act(ctx context.Context, vp *ValidatedExecutionPlan) error {
	if vp.Disposition != "approved" {
		return fmt.Errorf("act refused: plan not approved")
	}

	for _, step := range vp.Plan.Steps {
		fmt.Printf("Act: Executing step: %s\n", step)
	}
	return nil
}
