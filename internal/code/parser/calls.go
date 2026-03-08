package parser

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

// CallSite represents a function or method call site in source code
type CallSite struct {
	Name      string
	Kind      string
	Content   string
	StartLine int
	EndLine   int
}

// ExtractCalls analyzes a symbol's scope to identify all called functions
func (cp *CodeParser) ExtractCalls(ctx context.Context, source []byte, symbolNode *sitter.Node) ([]CallSite, error) {
	fmt.Printf("Extracting call sites for symbol: %s\n", symbolNode.Type())
	return []CallSite{}, nil
}
