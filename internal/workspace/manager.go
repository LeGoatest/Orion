package workspace

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"orion/internal/storage/sqlite"
)

type Manager struct {
	db   *sql.DB // Global DB
	dir  string
	mu   sync.RWMutex
	dbs  map[string]*sql.DB
}

func NewManager(db *sql.DB, d string) *Manager {
	return &Manager{
		db:  db,
		dir: d,
		dbs: make(map[string]*sql.DB),
	}
}

func (m *Manager) CreateWorkspace(id string) (string, error) {
	p := filepath.Join(m.dir, "workspaces", id)
	if err := os.MkdirAll(p, 0755); err != nil {
		return "", err
	}

	dbPath := filepath.Join(p, "workspace.db")
	_, err := sqlite.InitializeWorkspaceDB(dbPath)
	if err != nil {
		return "", err
	}

	// Register in global DB
	_, err = m.db.Exec("INSERT INTO workspaces (id, path) VALUES (?, ?) ON CONFLICT(id) DO NOTHING", id, p)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *Manager) GetWorkspaceDB(id string) (*sql.DB, error) {
	m.mu.RLock()
	if db, ok := m.dbs[id]; ok {
		m.mu.RUnlock()
		return db, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check
	if db, ok := m.dbs[id]; ok {
		return db, nil
	}

	p := filepath.Join(m.dir, "workspaces", id, "workspace.db")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace %s does not exist", id)
	}

	db, err := sqlite.InitializeWorkspaceDB(p)
	if err != nil {
		return nil, err
	}

	m.dbs[id] = db
	return db, nil
}
