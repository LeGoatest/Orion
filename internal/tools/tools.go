package tools

import (
	"context"
	"fmt"
	"os/exec"
)

type ShellTool struct{}

func (s *ShellTool) Name() string        { return "shell" }
func (s *ShellTool) Description() string { return "Executes shell commands" }
func (s *ShellTool) Safety() SafetyLevel { return LevelDangerous }
func (s *ShellTool) Execute(ctx context.Context, input string) (string, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

type ReadFileTool struct{}

func (r *ReadFileTool) Name() string        { return "read_file" }
func (r *ReadFileTool) Description() string { return "Reads a file" }
func (r *ReadFileTool) Safety() SafetyLevel { return LevelSafe }
func (r *ReadFileTool) Execute(ctx context.Context, input string) (string, error) {
	// Logic to read file
	return "file content", nil
}
