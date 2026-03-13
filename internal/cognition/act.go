package cognition

import (
	"context"
	"fmt"
)

func Act(ctx context.Context, plan *ExecutionPlan) error {
	for _, step := range plan.Steps {
		fmt.Printf("Act: %s\n", step)
	}
	return nil
}
