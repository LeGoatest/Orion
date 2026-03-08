package agents

import (
	"context"
	"fmt"
	"orion/internal/tools"
)

type ToolExecutorAgent struct {
	name         string
	capabilities []string
	registry     *tools.ToolRegistry
}

func NewToolExecutorAgent(registry *tools.ToolRegistry) *ToolExecutorAgent {
	return &ToolExecutorAgent{
		name:         "ToolExecutorAgent",
		capabilities: []string{"act_stage_execution", "tool_execution"},
		registry:     registry,
	}
}

func (a *ToolExecutorAgent) Name() string           { return a.name }
func (a *ToolExecutorAgent) Capabilities() []string { return a.capabilities }
func (a *ToolExecutorAgent) Priority() int          { return 2 }

func (a *ToolExecutorAgent) ExecuteTask(ctx context.Context, task interface{}) (interface{}, error) {
	fmt.Printf("Agent [%s]: Executing tool task\n", a.name)
	// Logic to execute tool from task description
	return "execution result", nil
}
