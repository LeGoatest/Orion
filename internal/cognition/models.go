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
	GoalID               string
	WorkspaceID          string
	NormalizedEvent      *NormalizedEvent
	ActiveSymbols        []*ent.CodeSymbol
	PatternMatches       []*ent.Pattern
	MemoryHits           []*ent.MemoryNode
	CapabilityCandidates []string
	ConstraintSet        []string
	ConfidenceSignals    map[string]float64
	OrientationSummary   string
	Timestamp            time.Time
}

type ExecutionPlan struct {
	GoalID string
	Steps  []string
	Tools  []string
}

type ValidatedExecutionPlan struct {
	Plan              *ExecutionPlan
	Disposition       string // approved, rejected, requires_revision
	ValidationReason  string
	ValidatedAt       time.Time
}

type OutcomeRecord struct {
	GoalID    string
	Success   bool
	Result    interface{}
	Timestamp time.Time
}
