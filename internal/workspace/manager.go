package workspace

import "database/sql"
import "path/filepath"
import "os"
import "orion/internal/storage/sqlite"

type Manager struct { db *sql.DB; dir string }
func NewManager(db *sql.DB, d string) *Manager { return &Manager{db: db, dir: d} }
func (m *Manager) CreateWorkspace(n string) (string, error) { p := filepath.Join(m.dir, "workspaces", n); os.MkdirAll(p, 0755); sqlite.InitializeWorkspaceDB(filepath.Join(p, "workspace.db")); return n, nil }
