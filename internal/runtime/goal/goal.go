package goal

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
	StagePlanning     GoalStage = "PLANNING"
	StageExecution    GoalStage = "EXECUTION"
	StageLearn        GoalStage = "LEARN"
	StageGarden       GoalStage = "GARDEN"
)

type Goal struct {
	ID            string
	Description   string
	CurrentStage  GoalStage
	AssignedAgent string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type GoalRuntime struct{}

func (gr *GoalRuntime) AdvanceStage(g *Goal, next GoalStage, agentName string) {
	fmt.Printf("Goal %s: %s -> %s (Agent: %s)\n", g.ID, g.CurrentStage, next, agentName)
	g.CurrentStage = next
	g.AssignedAgent = agentName
	g.UpdatedAt = time.Now()
}
