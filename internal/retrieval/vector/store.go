package vector

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// SearchResult contains the node and its similarity score
type SearchResult struct {
	ID    string
	Score float32
}

// VectorStore handles semantic vector search using sqlite_vec
type VectorStore struct {
	db *sql.DB
}

func NewVectorStore(db *sql.DB) *VectorStore {
	return &VectorStore{db: db}
}

// Search retrieves the top-k similar memories for a given embedding
func (vs *VectorStore) Search(ctx context.Context, embedding []float32, k int) ([]SearchResult, error) {
	embJSON, _ := json.Marshal(embedding)

	query := `SELECT id, distance FROM memory_embeddings WHERE embedding MATCH ? AND k = ?`
	rows, err := vs.db.QueryContext(ctx, query, string(embJSON), k)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var res SearchResult
		if err := rows.Scan(&res.ID, &res.Score); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

// StoreEmbedding adds a memory node's embedding to the vector table
func (vs *VectorStore) StoreEmbedding(ctx context.Context, id string, embedding []float32) error {
	embJSON, _ := json.Marshal(embedding)
	query := `INSERT OR REPLACE INTO memory_embeddings (id, embedding) VALUES (?, ?)`
	_, err := vs.db.ExecContext(ctx, query, id, string(embJSON))
	return err
}
