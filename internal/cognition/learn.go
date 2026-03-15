package cognition

import (
	"context"
	"orion/ent"
	"orion/ent/goal"
	"time"
)

func Learn(ctx context.Context, client *ent.Client, plan *ValidatedExecutionPlan) (*OutcomeRecord, error) {
	err := client.Goal.Update().
		Where(goal.ID(plan.Plan.GoalID)).
		SetStatus("completed").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &OutcomeRecord{
		GoalID:    plan.Plan.GoalID,
		Success:   true,
		Result:    "success",
		Timestamp: time.Now(),
	}, nil
}
