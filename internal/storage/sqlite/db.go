package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

func init() {
	// Register a custom sqlite3 driver that allows loading extensions
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		Extensions: []string{
			"sqlite_vec", // This assumes sqlite_vec is in the library path
		},
	})
}

// GlobalDBSchema defines the schema for the global orion.db
const GlobalDBSchema = `
CREATE TABLE IF NOT EXISTS workspaces (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	path TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS agents (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	status TEXT NOT NULL,
	last_active DATETIME
);

CREATE TABLE IF NOT EXISTS agent_status (
	agent_id TEXT PRIMARY KEY,
	status TEXT NOT NULL,
	last_heartbeat DATETIME,
	metadata TEXT,
	FOREIGN KEY (agent_id) REFERENCES agents(id)
);

CREATE TABLE IF NOT EXISTS system_config (
	key TEXT PRIMARY KEY,
	value TEXT
);

CREATE TABLE IF NOT EXISTS model_registry (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	provider TEXT NOT NULL,
	type TEXT NOT NULL, -- embedding, chat, etc
	config TEXT
);

CREATE TABLE IF NOT EXISTS event_log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_type TEXT NOT NULL,
	source TEXT NOT NULL,
	payload TEXT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

// WorkspaceDBSchema defines the schema for workspace-specific databases
const WorkspaceDBSchema = `
-- Goals and events
CREATE TABLE IF NOT EXISTS goals (
	id TEXT PRIMARY KEY,
	description TEXT NOT NULL,
	status TEXT NOT NULL,
	priority INTEGER DEFAULT 0,
	parent_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goal_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	goal_id TEXT NOT NULL,
	event_type TEXT NOT NULL,
	payload TEXT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (goal_id) REFERENCES goals(id)
);

-- Memory Graph
CREATE TABLE IF NOT EXISTS memory_nodes (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	content TEXT NOT NULL,
	importance REAL DEFAULT 0.0,
	usage_count INTEGER DEFAULT 0,
	archived BOOLEAN DEFAULT FALSE,
	metadata TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS memory_links (
	id TEXT PRIMARY KEY,
	source_id TEXT NOT NULL,
	target_id TEXT NOT NULL,
	relation_type TEXT NOT NULL,
	metadata TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (source_id) REFERENCES memory_nodes(id),
	FOREIGN KEY (target_id) REFERENCES memory_nodes(id)
);

CREATE TABLE IF NOT EXISTS patterns (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	pattern_data TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pattern_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	pattern_id TEXT NOT NULL,
	event_type TEXT NOT NULL,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (pattern_id) REFERENCES patterns(id)
);

-- Code Intelligence
CREATE TABLE IF NOT EXISTS code_symbols (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	kind TEXT NOT NULL,
	file_path TEXT NOT NULL,
	start_line INTEGER,
	end_line INTEGER,
	signature TEXT,
	metadata TEXT
);

CREATE TABLE IF NOT EXISTS call_graph_edges (
	id TEXT PRIMARY KEY,
	caller_id TEXT NOT NULL,
	callee_id TEXT NOT NULL,
	call_site TEXT,
	FOREIGN KEY (caller_id) REFERENCES code_symbols(id),
	FOREIGN KEY (callee_id) REFERENCES code_symbols(id)
);
`

// InitializeDB opens and initializes a SQLite database at the given path with the provided schema.
func InitializeDB(dbPath string, schema string) (*sql.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for DB: %w", err)
	}

	// Standard sqlite3 driver
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Try custom driver for extension support
	customDB, err := sql.Open("sqlite3_custom", dbPath)
	if err == nil {
		// Ping to ensure connection is valid and extensions loaded
		if err := customDB.Ping(); err == nil {
			db.Close()
			db = customDB
			fmt.Println("SQLite initialized with custom extension driver")
		} else {
			customDB.Close()
			fmt.Printf("Warning: custom sqlite driver ping failed: %v\n", err)
		}
	} else {
		fmt.Printf("Warning: failed to open with custom sqlite3 driver: %v\n", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// SetupVectorTable initializes the sqlite_vec virtual table for embeddings.
func SetupVectorTable(db *sql.DB) error {
	// Example sqlite_vec initialization
	_, err := db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS memory_embeddings USING vec0(
	    id TEXT PRIMARY KEY,
	    embedding FLOAT32[1536]
	)`)
	return err
}
