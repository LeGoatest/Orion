package indexer

import (
	"context"
	"fmt"
	"orion/internal/code/parser"
)

type CodeSymbol struct {
	ID        string
	Name      string
	Type      string
	FilePath  string
	StartLine int
	EndLine   int
}

type Indexer struct {
	parser *parser.CodeParser
}

func NewIndexer(p *parser.CodeParser) *Indexer {
	return &Indexer{parser: p}
}

func (idx *Indexer) IndexFile(ctx context.Context, filePath string, source []byte) error {
	tree, err := idx.parser.Parse(ctx, source)
	if err != nil {
		return err
	}

	fmt.Printf("Indexed file: %s, root type: %s\n", filePath, tree.RootNode().Type())
	// Logic to traverse tree and extract symbols
	return nil
}
