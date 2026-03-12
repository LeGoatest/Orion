package cognition

import (
	"orion/internal/symbols"
	"time"
)

type SituationalModel struct {
	GoalID           string
	GoalContext      string
	WorkspaceContext string
	CodeContext      []symbols.Symbol
	PatternContext   []string
	RetrievalContext []string
	Capabilities     []string
	GovernanceRules  []string
	Timestamp        time.Time
}

type DecisionModel struct {
	SituationalModel *SituationalModel
	PlanSteps        []string
	SelectedTools    []string
	StrategyChoice   string
	Timestamp        time.Time
}

type OutcomeModel struct {
	GoalID         string
	DecisionModel  *DecisionModel
	ActionResult   interface{}
	Success        bool
	OutcomeSummary string
	Timestamp      time.Time
}
