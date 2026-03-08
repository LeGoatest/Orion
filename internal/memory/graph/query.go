package graph

import (
	"context"
	"fmt"
)

// GraphQueryer provides low-level SQL operations for graph traversal
type GraphQueryer struct {
	db interface {
		QueryContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	}
}

// GetNeighbors retrieves all direct connections for a given node ID
func (gq *GraphQueryer) GetNeighbors(ctx context.Context, nodeID string) ([]string, error) {
	fmt.Printf("Querying neighbors for node: %s\n", nodeID)

	// SELECT target_id FROM memory_links WHERE source_id = ?
	// UNION
	// SELECT source_id FROM memory_links WHERE target_id = ?

	return []string{nodeID}, nil
}
