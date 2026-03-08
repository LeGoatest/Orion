package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeGlobalDB(p string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", p+"?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	if err := initGlobalSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InitializeWorkspaceDB(p string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", p+"?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	if err := initWorkspaceSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initGlobalSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT PRIMARY KEY,
			path TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed to exec query: %w", err)
		}
	}
	return nil
}

func initWorkspaceSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS goals (
			id TEXT PRIMARY KEY,
			description TEXT NOT NULL,
			current_stage TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			goal_id TEXT NOT NULL,
			stage TEXT NOT NULL,
			assigned_agent TEXT NOT NULL,
			status TEXT NOT NULL,
			retry_count INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			completed_at DATETIME,
			FOREIGN KEY(goal_id) REFERENCES goals(id)
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed to exec query: %w", err)
		}
	}
	return nil
}
