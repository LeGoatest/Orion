package agents

import (
	"context"
	"fmt"
)

type BaseAgent struct {
	AgentName string
	AgentType string
}

func (b *BaseAgent) Name() string { return b.AgentName }
func (b *BaseAgent) Type() string { return b.AgentType }

func (b *BaseAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return fmt.Sprintf("%s executed", b.AgentName), nil
}

func NewConversationAgent() Agent { return &BaseAgent{"conversation_agent", "communication"} }
func NewPlannerAgent() Agent      { return &BaseAgent{"planner_agent", "planning"} }
func NewAnalysisAgent() Agent     { return &BaseAgent{"analysis_agent", "analysis"} }
func NewCodeIndexerAgent() Agent  { return &BaseAgent{"code_indexer_agent", "indexing"} }
func NewMemoryGardenerAgent() Agent { return &BaseAgent{"memory_gardener_agent", "maintenance"} }
func NewPatternDetectorAgent() Agent { return &BaseAgent{"pattern_detector_agent", "analysis"} }
