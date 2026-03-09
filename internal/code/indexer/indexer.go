package indexer

import (
	"context"
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	sitter_go "github.com/smacker/go-tree-sitter/go"
	"orion/internal/storage/sqlite"
	"os"
	// sitter_py "github.com/smacker/go-tree-sitter/python"
	// sitter_ts "github.com/smacker/go-tree-sitter/typescript"
	// sitter_js "github.com/smacker/go-tree-sitter/javascript"
)

type Indexer struct {
	db *sqlite.DB
}

func NewIndexer(db *sqlite.DB) *Indexer {
	return &Indexer{db: db}
}

func (idx *Indexer) IndexFile(ctx context.Context, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lang := idx.getLanguage(filePath)
	if lang == nil {
		return fmt.Errorf("unsupported language for %s", filePath)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, content)
	if err != nil {
		return err
	}

	root := tree.RootNode()
	// Extract symbols from the root node and store in sqlite
	fmt.Printf("Indexed file %s: %d nodes\n", filePath, root.ChildCount())

	return nil
}

func (idx *Indexer) getLanguage(filePath string) *sitter.Language {
	switch filepath.Ext(filePath) {
	case ".go":
		return sitter_go.GetLanguage()
	// Other languages can be added here once their grammars are imported
	default:
		return nil
	}
}
