package tools

import (
	"context"
	"fmt"
)

type SafetyLevel string

const (
	LevelSafe       SafetyLevel = "safe"
	LevelRestricted SafetyLevel = "restricted"
	LevelDangerous  SafetyLevel = "dangerous"
)

type Tool interface {
	Name() string
	Description() string
	Safety() SafetyLevel
	Execute(ctx context.Context, input string) (string, error)
}

type ToolRegistry struct {
	tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]Tool)}
}

func (tr *ToolRegistry) Register(tool Tool) {
	tr.tools[tool.Name()] = tool
}

func (tr *ToolRegistry) ExecuteTool(ctx context.Context, name, input string) (string, error) {
	tool, ok := tr.tools[name]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	if tool.Safety() == LevelDangerous {
		fmt.Printf("Warning: Executing dangerous tool %s\n", name)
		// Policy approval logic
	}

	return tool.Execute(ctx, input)
}
