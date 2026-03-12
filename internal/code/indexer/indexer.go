package indexer

import (
	"fmt"
)

type Indexer struct {
	Languages []string
}

func NewIndexer() *Indexer {
	return &Indexer{
		Languages: []string{"Go", "TypeScript", "JavaScript", "Markdown"},
	}
}

func (i *Indexer) IndexFile(path string) error {
	fmt.Printf("Indexing file: %s\n", path)
	// Skeleton for Tree-Sitter indexing
	return nil
}
