package goal

import (
	"context"
	"fmt"
)

type ObserveJob struct {
	GoalID string
}

func (j *ObserveJob) ID() string { return j.GoalID + "-observe" }
func (j *ObserveJob) Type() string { return "ObserveJob" }
func (j *ObserveJob) Execute(ctx context.Context) error {
	fmt.Printf("Executing ObserveJob for Goal %s\n", j.GoalID)
	// Logic to capture raw intent and create memory node
	return nil
}

type OrientJob struct {
	GoalID string
}

func (j *OrientJob) ID() string { return j.GoalID + "-orient" }
func (j *OrientJob) Type() string { return "OrientJob" }
func (j *OrientJob) Execute(ctx context.Context) error {
	fmt.Printf("Executing OrientJob for Goal %s\n", j.GoalID)
	// Logic for hybrid retrieval
	return nil
}
