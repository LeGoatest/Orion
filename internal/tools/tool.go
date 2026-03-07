package tools

import (
	"context"
	"fmt"
)

// BaseTool provides common functionality for tools
type BaseTool struct {
	name        string
	description string
}

func (t *BaseTool) Name() string {
	return t.name
}

func (t *BaseTool) Description() string {
	return t.description
}

// ShellTool example
type ShellTool struct {
	BaseTool
}

func NewShellTool() *ShellTool {
	return &ShellTool{
		BaseTool: BaseTool{
			name:        "shell",
			description: "Executes shell commands",
		},
	}
}

func (t *ShellTool) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	// Shell execution logic
	return fmt.Sprintf("Executed shell: %v", input), nil
}
