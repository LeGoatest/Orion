package workspace

import (
	"database/sql"
	"fmt"
	"orion/internal/storage/sqlite"
	"path/filepath"
	"sync"
)

type Workspace struct {
	ID   string
	Name string
	DB   *sql.DB
}

type Manager struct {
	dataDir    string
	globalDB   *sql.DB
	workspaces map[string]*Workspace
	mu         sync.RWMutex
}

func NewManager(globalDB *sql.DB, dataDir string) *Manager {
	return &Manager{
		dataDir:    dataDir,
		globalDB:   globalDB,
		workspaces: make(map[string]*Workspace),
	}
}

func (m *Manager) Start() {
	rows, err := m.globalDB.Query("SELECT id, name, path FROM workspaces")
	if err != nil {
		fmt.Printf("Error loading workspaces: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, path string
		if err := rows.Scan(&id, &name, &path); err != nil {
			fmt.Printf("Error scanning workspace: %v\n", err)
			continue
		}

		db, err := sqlite.OpenWorkspace(path)
		if err != nil {
			fmt.Printf("Error opening workspace database for %s: %v\n", id, err)
			continue
		}

		m.mu.Lock()
		m.workspaces[id] = &Workspace{
			ID:   id,
			Name: name,
			DB:   db,
		}
		m.mu.Unlock()
	}
}

func (m *Manager) CreateWorkspace(id, name string) (*Workspace, error) {
	workspacePath := filepath.Join(m.dataDir, "workspaces", id)
	db, err := sqlite.OpenWorkspace(workspacePath)
	if err != nil {
		return nil, err
	}

	_, err = m.globalDB.Exec("INSERT INTO workspaces (id, name, path) VALUES (?, ?, ?)", id, name, workspacePath)
	if err != nil {
		db.Close()
		return nil, err
	}

	w := &Workspace{
		ID:   id,
		Name: name,
		DB:   db,
	}

	m.mu.Lock()
	m.workspaces[id] = w
	m.mu.Unlock()

	return w, nil
}
