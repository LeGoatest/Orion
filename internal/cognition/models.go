package cognition

import (
	"orion/ent"
	"time"
)

type NormalizedEvent struct {
	ID        string
	GoalID    string
	Type      string
	Payload   map[string]interface{}
	Timestamp time.Time
}

type SituationalModel struct {
	Goal             *ent.Goal
	WorkspaceContext string
	Symbols          []*ent.CodeSymbol
	Memories         []*ent.MemoryNode
	Patterns         []*ent.Pattern
	Timestamp        time.Time
}

type ExecutionPlan struct {
	GoalID string
	Steps  []string
	Tools  []string
}

type OutcomeRecord struct {
	GoalID    string
	Success   bool
	Result    interface{}
	Timestamp time.Time
}
