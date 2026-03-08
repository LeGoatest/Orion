package goal

import (
	"fmt"
	"time"
)

type GoalState string

const (
	StateNew       GoalState = "NEW"
	StateObserve   GoalState = "OBSERVE"
	StateOrient    GoalState = "ORIENT"
	StateDecide    GoalState = "DECIDE"
	StateAct       GoalState = "ACT"
	StateLearn     GoalState = "LEARN"
	StateGarden    GoalState = "GARDEN"
	StateWaiting   GoalState = "WAITING"
	StateCompleted GoalState = "COMPLETED"
	StateFailed    GoalState = "FAILED"
	StateArchived  GoalState = "ARCHIVED"
)

type Goal struct {
	ID          string
	Description string
	State       GoalState
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GoalRuntime struct{}

func NewGoalRuntime() *GoalRuntime {
	return &GoalRuntime{}
}

func (gr *GoalRuntime) Transition(goal *Goal, newState GoalState) {
	fmt.Printf("Goal %s: %s -> %s\n", goal.ID, goal.State, newState)
	goal.State = newState
	goal.UpdatedAt = time.Now()
}
