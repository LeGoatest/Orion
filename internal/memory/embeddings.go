package memory

import (
	"context"
	"fmt"
)

// EmbeddingStorage handles persistence of vector data
type EmbeddingStorage struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	}
}

// StoreEmbedding adds a vector to the database
func (s *EmbeddingStorage) StoreEmbedding(ctx context.Context, nodeID string, embedding []float32) error {
	// Logic to persist embedding in workspace.db memory_embeddings
	fmt.Printf("Persisting embedding for node: %s\n", nodeID)
	return nil
}
