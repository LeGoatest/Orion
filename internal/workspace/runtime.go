package workspace

import (
	"context"
	"fmt"
	"sync"
)

// WorkspaceRuntime manages the execution of a single workspace
type WorkspaceRuntime struct {
	Workspace *Workspace
	mu        sync.RWMutex
	running   bool
}

// NewWorkspaceRuntime creates a new runtime for a workspace
func NewWorkspaceRuntime(ws *Workspace) *WorkspaceRuntime {
	return &WorkspaceRuntime{
		Workspace: ws,
	}
}

// Start initiates workspace-specific services and background tasks
func (wr *WorkspaceRuntime) Start(ctx context.Context) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if wr.running {
		return fmt.Errorf("workspace %s is already running", wr.Workspace.ID)
	}

	fmt.Printf("Starting runtime for workspace: %s (%s)\n", wr.Workspace.Name, wr.Workspace.ID)

	// Workspace-specific initialization (e.g. memory gardener for this workspace)
	wr.running = true
	return nil
}

// Stop gracefully shuts down workspace-specific services
func (wr *WorkspaceRuntime) Stop(ctx context.Context) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if !wr.running {
		return nil
	}

	fmt.Printf("Stopping runtime for workspace: %s\n", wr.Workspace.ID)
	wr.running = false
	return nil
}
