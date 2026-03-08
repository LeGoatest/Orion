package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"orion/internal/storage/sqlite"
	"github.com/google/uuid"
)

// Workspace represents an isolated cognitive environment
type Workspace struct {
	ID   string
	Name string
	Path string
	DB   *sql.DB
}

// Manager handles workspace lifecycle and registry
type Manager struct {
	mu         sync.RWMutex
	globalDB   *sql.DB
	dataDir    string
	workspaces map[string]*Workspace
}

// NewManager creates a new Manager
func NewManager(globalDB *sql.DB, dataDir string) *Manager {
	return &Manager{
		globalDB:   globalDB,
		dataDir:    dataDir,
		workspaces: make(map[string]*Workspace),
	}
}

// CreateWorkspace initializes a new workspace
func (m *Manager) CreateWorkspace(ctx context.Context, name string) (*Workspace, error) {
	id := uuid.New().String()
	wsPath := filepath.Join(m.dataDir, "workspaces", id)
	dbPath := filepath.Join(wsPath, "workspace.db")

	// Ensure workspace directories exist
	dirs := []string{
		wsPath,
		filepath.Join(wsPath, "index"),
		filepath.Join(wsPath, "symbols"),
		filepath.Join(wsPath, "artifacts"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create workspace directory %s: %w", dir, err)
		}
	}

	// Initialize workspace DB
	db, err := sqlite.InitializeDB(dbPath, sqlite.WorkspaceDBSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize workspace DB: %w", err)
	}

	// Register in global DB
	query := `INSERT INTO workspaces (id, name, path) VALUES (?, ?, ?)`
	_, err = m.globalDB.ExecContext(ctx, query, id, name, wsPath)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to register workspace in global DB: %w", err)
	}

	ws := &Workspace{
		ID:   id,
		Name: name,
		Path: wsPath,
		DB:   db,
	}

	m.mu.Lock()
	m.workspaces[id] = ws
	m.mu.Unlock()

	return ws, nil
}

// LoadWorkspaces loads existing workspaces from the global DB
func (m *Manager) LoadWorkspaces(ctx context.Context) error {
	query := `SELECT id, name, path FROM workspaces`
	rows, err := m.globalDB.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query workspaces: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ws Workspace
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.Path); err != nil {
			return fmt.Errorf("failed to scan workspace: %w", err)
		}

		dbPath := filepath.Join(ws.Path, "workspace.db")
		db, err := sqlite.InitializeDB(dbPath, sqlite.WorkspaceDBSchema)
		if err != nil {
			fmt.Printf("Warning: failed to load workspace DB for %s: %v\n", ws.ID, err)
			continue
		}
		ws.DB = db

		m.mu.Lock()
		m.workspaces[ws.ID] = &ws
		m.mu.Unlock()
	}

	return nil
}

// GetWorkspace returns a workspace by ID
func (m *Manager) GetWorkspace(id string) (*Workspace, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ws, ok := m.workspaces[id]
	return ws, ok
}

// ListWorkspaces returns all loaded workspaces
func (m *Manager) ListWorkspaces() []*Workspace {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]*Workspace, 0, len(m.workspaces))
	for _, ws := range m.workspaces {
		list = append(list, ws)
	}
	return list
}
