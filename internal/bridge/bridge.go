package bridge

import (
	"context"
	"fmt"
	"orion/internal/runtime"
	"orion/internal/runtime/goal"
	"orion/internal/types"
	"time"
)

// Bridge connects the Wails frontend to the Orion runtime kernel
type Bridge struct {
	ctx    context.Context
	kernel *runtime.Kernel
}

func NewBridge(kernel *runtime.Kernel) *Bridge {
	return &Bridge{
		kernel: kernel,
	}
}

// Startup is called by Wails when the application starts
func (b *Bridge) Startup(ctx context.Context) {
	b.ctx = ctx
}

// SubmitGoal creates a new goal in the current workspace
func (b *Bridge) SubmitGoal(workspaceID, description string) string {
	fmt.Printf("UI Action: SubmitGoal - %s\n", description)
	// In a real implementation, this would call kernel.GetWorkspace(workspaceID).SubmitGoal(description)
	return "Goal submitted successfully"
}

// GetGoals retrieves all goals for a workspace
func (b *Bridge) GetGoals(workspaceID string) []goal.Goal {
	fmt.Printf("UI Action: GetGoals for %s\n", workspaceID)
	return []goal.Goal{}
}

// GetMemoryNodes retrieves nodes for the memory explorer
func (b *Bridge) GetMemoryNodes(workspaceID string) []interface{} {
	fmt.Printf("UI Action: GetMemoryNodes for %s\n", workspaceID)
	return []interface{}{}
}

// GetCodeIntelligence retrieves symbols and call graph for a workspace
func (b *Bridge) GetCodeIntelligence(workspaceID string) map[string]interface{} {
	fmt.Printf("UI Action: GetCodeIntelligence for %s\n", workspaceID)
	return make(map[string]interface{})
}

// ListWorkspaces returns the available workspaces from the global registry
func (b *Bridge) ListWorkspaces() []string {
	return []string{"alpha", "beta"}
}

// GetAgentStatus returns the current state of all registered agents
func (b *Bridge) GetAgentStatus() []map[string]string {
	return []map[string]string{
		{"name": "ConversationAgent", "status": "idle"},
		{"name": "PlannerAgent", "status": "active"},
	}
}
