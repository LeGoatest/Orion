package runtime

import (
	"orion/internal/types"
	"time"
)

type GoalStage string

const (
	StageObserve      GoalStage = "OBSERVE"
	StageOrient       GoalStage = "ORIENT"
	StageDecide       GoalStage = "DECIDE"
	StageAct          GoalStage = "ACT"
	StageLearn        GoalStage = "LEARN"
	StageCompleted    GoalStage = "COMPLETED"
	StageFailed       GoalStage = "FAILED"
)

type Goal struct {
	ID            string
	Description   string
	CurrentStage  GoalStage
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type BaseAgent struct {
	EventBus *types.EventBus
}
