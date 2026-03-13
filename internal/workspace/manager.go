package workspace

import (
	"context"
	"fmt"
	"orion/ent"
	"sync"
)

type Workspace struct {
	ID   string
	Name string
	DB   *ent.Client
}

type Manager struct {
	dataDir    string
	client     *ent.Client
	workspaces map[string]*Workspace
	mu         sync.RWMutex
}

func NewManager(client *ent.Client, dataDir string) *Manager {
	return &Manager{
		dataDir:    dataDir,
		client:     client,
		workspaces: make(map[string]*Workspace),
	}
}

func (m *Manager) Start(ctx context.Context) {
	ws, err := m.client.Workspace.Query().All(ctx)
	if err != nil {
		fmt.Printf("Error loading workspaces: %v\n", err)
		return
	}

	for _, w := range ws {
		fmt.Printf("Loading workspace: %s at %s\n", w.Name, w.Path)
	}
}

func (m *Manager) CreateWorkspace(ctx context.Context, name string) (*WorkspaceRuntime, error) {
	path := fmt.Sprintf("%s/workspaces/%s.db", m.dataDir, name)

	_, err := m.client.Workspace.Create().
		SetName(name).
		SetPath(path).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return NewWorkspaceRuntime(name, path)
}
