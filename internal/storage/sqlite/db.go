package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// GlobalDBSchema defines the schema for the global orion.db
const GlobalDBSchema = `
CREATE TABLE IF NOT EXISTS workspaces (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	path TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

// WorkspaceDBSchema defines the schema for workspace-specific databases
const WorkspaceDBSchema = `
CREATE TABLE IF NOT EXISTS symbols (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	file_path TEXT NOT NULL,
	start_line INTEGER,
	end_line INTEGER,
	workspace_id TEXT,
	metadata TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_symbols_name ON symbols(name);
CREATE INDEX IF NOT EXISTS idx_symbols_type ON symbols(type);
CREATE INDEX IF NOT EXISTS idx_symbols_file_path ON symbols(file_path);

CREATE TABLE IF NOT EXISTS patterns (
	id TEXT PRIMARY KEY,
	trigger TEXT NOT NULL,
	solution_steps TEXT NOT NULL,
	confidence REAL,
	usage_count INTEGER,
	state TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

// InitializeGlobalDB initializes the global orion.db
func InitializeGlobalDB(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open global db: %w", err)
	}

	if _, err := db.Exec(GlobalDBSchema); err != nil {
		return nil, fmt.Errorf("failed to initialize global schema: %w", err)
	}

	return db, nil
}

// InitializeWorkspaceDB initializes a workspace-specific database
func InitializeWorkspaceDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open workspace db: %w", err)
	}

	if _, err := db.Exec(WorkspaceDBSchema); err != nil {
		return nil, fmt.Errorf("failed to initialize workspace schema: %w", err)
	}

	return db, nil
}
