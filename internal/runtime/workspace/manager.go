package runtime

import (
	"database/sql"
	"fmt"
)

type WorkspaceManager struct {
	db      *sql.DB
	dataDir string
}

func NewWorkspaceManager(db *sql.DB, dataDir string) *WorkspaceManager {
	return &WorkspaceManager{
		db:      db,
		dataDir: dataDir,
	}
}

func (wm *WorkspaceManager) CreateWorkspace(name string) (string, error) {
	fmt.Printf("Creating workspace: %s\n", name)
	// Path creation, registry entry
	return "ws-id", nil
}
