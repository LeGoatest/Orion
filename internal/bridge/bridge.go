package bridge

import (
	"context"
	"fmt"
)

// Bridge connects the UI to the runtime kernel
type Bridge struct {
	ctx context.Context
}

func NewBridge() *Bridge {
	return &Bridge{}
}

func (b *Bridge) SetContext(ctx context.Context) {
	b.ctx = ctx
}

func (b *Bridge) HandleChat(msg string) string {
	fmt.Printf("UI Request: %s\n", msg)
	// Process via runtime kernel
	return "Acknowledged: " + msg
}

func (b *Bridge) ListWorkspaces() []string {
	return []string{"Workspace 1", "Workspace 2"}
}
