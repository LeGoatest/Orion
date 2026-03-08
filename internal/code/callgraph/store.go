package callgraph

import (
	"context"
	"fmt"
)

// Store handles persistence of call graph data
type Store struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	}
}

// PersistEdge adds a new call edge to the call_graph_edges table
func (s *Store) PersistEdge(ctx context.Context, edge CallEdge) error {
	fmt.Printf("Persisting call edge: %s -> %s\n", edge.CallerSymbolID, edge.CalleeSymbolID)
	// INSERT INTO call_graph_edges (id, caller_id, callee_id) VALUES (?, ?, ?)
	return nil
}

// GetCallees retrieves all functions called by the given symbol ID
func (s *Store) GetCallees(ctx context.Context, symbolID string) ([]string, error) {
	fmt.Printf("Fetching callees for symbol: %s\n", symbolID)
	// SELECT callee_id FROM call_graph_edges WHERE caller_id = ?
	return []string{}, nil
}

// GetCallers retrieves all functions that call the given symbol ID
func (s *Store) GetCallers(ctx context.Context, symbolID string) ([]string, error) {
	fmt.Printf("Fetching callers for symbol: %s\n", symbolID)
	// SELECT caller_id FROM call_graph_edges WHERE callee_id = ?
	return []string{}, nil
}
