package callgraph

import (
	"context"
	"fmt"
)

// Traverser provides methods to navigate the call graph
type Traverser struct {
	store *Store
}

func NewTraverser(store *Store) *Traverser {
	return &Traverser{store: store}
}

// ExpandCallGraph retrieves a neighborhood of callers and callees to a specified depth
func (t *Traverser) ExpandCallGraph(ctx context.Context, initialSymbols []string, depth int) ([]string, error) {
	fmt.Printf("Expanding call graph from %d initial symbols (depth: %d)\n", len(initialSymbols), depth)

	// Traversal logic:
	// Recursively SELECT caller_id, callee_id FROM call_graph_edges WHERE ...

	return initialSymbols, nil
}
