package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

func OpenGlobal(dataDir string) (*sql.DB, error) {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, err
		}
	}
	dbPath := filepath.Join(dataDir, "orion.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `CREATE TABLE IF NOT EXISTS workspaces (id TEXT PRIMARY KEY, name TEXT NOT NULL, path TEXT NOT NULL, created_at INTEGER DEFAULT (strftime('%s', 'now')));
	CREATE TABLE IF NOT EXISTS agents (id TEXT PRIMARY KEY, type TEXT NOT NULL, capabilities TEXT, status TEXT DEFAULT 'active');`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create global schema: %v", err)
	}
	return db, nil
}

func OpenWorkspace(workspacePath string) (*sql.DB, error) {
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		if err := os.MkdirAll(workspacePath, 0755); err != nil {
			return nil, err
		}
	}
	dbPath := filepath.Join(workspacePath, "workspace.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `CREATE TABLE IF NOT EXISTS goals (id TEXT PRIMARY KEY, description TEXT NOT NULL, current_stage TEXT NOT NULL DEFAULT 'OBSERVE', status TEXT DEFAULT 'pending', assigned_agent TEXT, created_at INTEGER DEFAULT (strftime('%s', 'now')), updated_at INTEGER DEFAULT (strftime('%s', 'now')));
	CREATE TABLE IF NOT EXISTS jobs (id TEXT PRIMARY KEY, goal_id TEXT NOT NULL, stage TEXT NOT NULL, assigned_agent TEXT NOT NULL, status TEXT NOT NULL, retry_count INTEGER DEFAULT 0, created_at INTEGER DEFAULT (strftime('%s', 'now')), finished_at INTEGER);`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create workspace schema: %v", err)
	}
	return db, nil
}
