package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"orion/internal/code/parser"
	"github.com/google/uuid"
)

// Indexer handles parsing and storing code intelligence data
type Indexer struct {
	db     *sql.DB
	parser parser.Parser
}

// NewIndexer creates a new Indexer
func NewIndexer(db *sql.DB, p parser.Parser) *Indexer {
	return &Indexer{
		db:     db,
		parser: p,
	}
}

// IndexFile parses a single file and stores symbols and call graph edges
func (i *Indexer) IndexFile(ctx context.Context, filePath string, lang parser.Language) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	nodes, err := i.parser.Parse(ctx, lang, content)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	for _, node := range nodes {
		symbolID := uuid.New().String()
		query := `INSERT INTO code_symbols (id, name, kind, file_path, start_line, end_line, signature)
                  VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err := i.db.ExecContext(ctx, query, symbolID, node.Name, node.Kind, filePath, node.StartLine, node.EndLine, node.Signature)
		if err != nil {
			fmt.Printf("Warning: failed to insert symbol %s: %v\n", node.Name, err)
			continue
		}

		// Simplified: record function calls to build call graph (not fully implemented)
		// This would involve additional parsing to find call sites within the current node.
	}

	return nil
}

// IndexProject crawls a directory and indexes all supported code files
func (i *Indexer) IndexProject(ctx context.Context, projectRoot string) error {
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		lang, supported := getLanguageFromExt(filepath.Ext(path))
		if supported {
			fmt.Printf("Indexing: %s\n", path)
			if err := i.IndexFile(ctx, path, lang); err != nil {
				fmt.Printf("Error indexing %s: %v\n", path, err)
			}
		}
		return nil
	})

	return err
}

func getLanguageFromExt(ext string) (parser.Language, bool) {
	switch ext {
	case ".go":
		return parser.Go, true
	case ".py":
		return parser.Python, true
	case ".js", ".jsx":
		return parser.JavaScript, true
	default:
		return "", false
	}
}
