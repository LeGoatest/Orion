package tools

import (
	"context"
	"sync"
)

// Tool defines the interface for an Orion tool
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

// Registry maintains a set of registered tools
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name()] = tool
}

// GetTool retrieves a tool by name
func (r *Registry) GetTool(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

// ListTools returns all registered tools
func (r *Registry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		list = append(list, tool)
	}
	return list
}
