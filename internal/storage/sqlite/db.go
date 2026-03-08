package sqlite

import "database/sql"
import "path/filepath"
import "os"
import _ "github.com/mattn/go-sqlite3"

func InitializeGlobalDB(p string) (*sql.DB, error) { os.MkdirAll(filepath.Dir(p), 0755); return sql.Open("sqlite3", p) }
func InitializeWorkspaceDB(p string) (*sql.DB, error) { os.MkdirAll(filepath.Dir(p), 0755); return sql.Open("sqlite3", p) }
