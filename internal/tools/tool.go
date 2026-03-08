package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// ShellTool: Executes system commands
type ShellTool struct {
	BaseTool
}

func NewShellTool() *ShellTool {
	return &ShellTool{
		BaseTool: BaseTool{
			name:        "shell",
			description: "Executes shell commands on the local machine",
		},
	}
}

func (t *ShellTool) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	cmdStr, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("invalid input: shell tool expects a command string")
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("command execution failed: %w", err)
	}

	return string(out), nil
}

// FileSystemTool: Handles file read/write operations
type FileSystemTool struct {
	BaseTool
}

func NewFileSystemTool() *FileSystemTool {
	return &FileSystemTool{
		BaseTool: BaseTool{
			name:        "filesystem",
			description: "Manages file system operations (read, write, list)",
		},
	}
}

func (t *FileSystemTool) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	params, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input: filesystem tool expects parameters as a map")
	}

	op, _ := params["operation"].(string)
	path, _ := params["path"].(string)

	switch op {
	case "read":
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read failed: %w", err)
		}
		return string(data), nil
	case "write":
		content, _ := params["content"].(string)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return nil, fmt.Errorf("write failed: %w", err)
		}
		return "success", nil
	case "list":
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("list failed: %w", err)
		}
		var list []string
		for _, e := range entries {
			list = append(list, e.Name())
		}
		return list, nil
	default:
		return nil, fmt.Errorf("unsupported operation: %s", op)
	}
}
