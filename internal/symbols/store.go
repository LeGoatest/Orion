package symbols

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Store struct {
	Db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Db: db}
}

// InsertSymbol persists a code symbol to the workspace database
func (s *Store) InsertSymbol(ctx context.Context, sym Symbol) error {
	id := uuid.New().String()
	query := `INSERT INTO symbols (id, name, type, file_path, start_line, end_line, workspace_id, metadata)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.Db.ExecContext(ctx, query, id, sym.Name, sym.Type, sym.FilePath, sym.StartLine, sym.EndLine, sym.WorkspaceID, sym.Metadata)
	return err
}

// FindByName performs an exact name search for symbols
func (s *Store) FindByName(ctx context.Context, name string) ([]Symbol, error) {
	query := `SELECT id, name, type, file_path, start_line, end_line, workspace_id, metadata FROM symbols WHERE name = ?`
	rows, err := s.Db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var syms []Symbol
	for rows.Next() {
		var sym Symbol
		err := rows.Scan(&sym.ID, &sym.Name, &sym.Type, &sym.FilePath, &sym.StartLine, &sym.EndLine, &sym.WorkspaceID, &sym.Metadata)
		if err != nil {
			return nil, err
		}
		syms = append(syms, sym)
	}
	return syms, nil
}

// FindFuzzy performs a pattern-based search for symbols
func (s *Store) FindFuzzy(ctx context.Context, query string) ([]Symbol, error) {
	sqlQuery := `SELECT id, name, type, file_path, start_line, end_line, workspace_id, metadata FROM symbols WHERE name LIKE ?`
	rows, err := s.Db.QueryContext(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var syms []Symbol
	for rows.Next() {
		var sym Symbol
		err := rows.Scan(&sym.ID, &sym.Name, &sym.Type, &sym.FilePath, &sym.StartLine, &sym.EndLine, &sym.WorkspaceID, &sym.Metadata)
		if err != nil {
			return nil, err
		}
		syms = append(syms, sym)
	}
	return syms, nil
}
