package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	path string
}

func OpenGlobal(dataDir string) (*DB, error) {
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

	schema := `
	CREATE TABLE IF NOT EXISTS workspaces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		capabilities TEXT,
		status TEXT DEFAULT 'active'
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create global schema: %v", err)
	}

	return &DB{DB: db, path: dbPath}, nil
}

func OpenWorkspace(workspacePath string) (*DB, error) {
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

	// Load sqlite_vec extension (CGO required)
	// _, err = db.Exec("SELECT load_extension('sqlite_vec');")

	schema := `
	CREATE TABLE IF NOT EXISTS goals (
		id TEXT PRIMARY KEY,
		description TEXT NOT NULL,
		parent_id TEXT,
		status TEXT DEFAULT 'pending',
		created_at INTEGER DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS goal_events (
		id TEXT PRIMARY KEY,
		goal_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		payload TEXT,
		timestamp INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS memory_nodes (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		importance REAL DEFAULT 0.5,
		usage_count INTEGER DEFAULT 0,
		archived BOOLEAN DEFAULT 0,
		created_at INTEGER DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS memory_links (
		id TEXT PRIMARY KEY,
		source_id TEXT NOT NULL,
		target_id TEXT NOT NULL,
		relationship TEXT NOT NULL,
		weight REAL DEFAULT 1.0,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS memory_embeddings (
		node_id TEXT PRIMARY KEY,
		embedding BLOB NOT NULL
	);
	CREATE TABLE IF NOT EXISTS patterns (
		id TEXT PRIMARY KEY,
		description TEXT NOT NULL,
		confidence REAL,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS pattern_events (
		id TEXT PRIMARY KEY,
		pattern_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		timestamp INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS pattern_embeddings (
		pattern_id TEXT PRIMARY KEY,
		embedding BLOB NOT NULL
	);
	CREATE TABLE IF NOT EXISTS code_symbols (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		file_path TEXT NOT NULL,
		line_start INTEGER,
		line_end INTEGER,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS symbol_embeddings (
		symbol_id TEXT PRIMARY KEY,
		embedding BLOB NOT NULL
	);
	CREATE TABLE IF NOT EXISTS symbol_references (
		id TEXT PRIMARY KEY,
		symbol_id TEXT NOT NULL,
		file_path TEXT NOT NULL,
		line INTEGER,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	CREATE TABLE IF NOT EXISTS call_graph_edges (
		id TEXT PRIMARY KEY,
		caller_id TEXT NOT NULL,
		callee_id TEXT NOT NULL,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create workspace schema: %v", err)
	}

	return &DB{DB: db, path: dbPath}, nil
}
