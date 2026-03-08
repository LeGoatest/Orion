package runtime

import (
	"fmt"
	"time"
)

type GoalStage string

const (
	StageObserve      GoalStage = "OBSERVE"
	StageSymbolLookup GoalStage = "SYMBOL_LOOKUP"
	StagePatternMatch GoalStage = "PATTERN_MATCH"
	StageRetrieval    GoalStage = "RETRIEVAL"
	StagePlan         GoalStage = "PLAN"
	StageAct          GoalStage = "ACT"
	StageLearn        GoalStage = "LEARN"
	StageGarden       GoalStage = "GARDEN"
	StageCompleted    GoalStage = "COMPLETED"
	StageFailed       GoalStage = "FAILED"
)

type Goal struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	CurrentStage  GoalStage `json:"current_stage"`
	AssignedAgent string    `json:"assigned_agent"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (g *Goal) AdvanceStage(next GoalStage) {
	fmt.Printf("Goal %s: %s -> %s\n", g.ID, g.CurrentStage, next)
	g.CurrentStage = next
	g.UpdatedAt = time.Now()
}
