package memory

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MemoryNode represents an atom of knowledge in Orion
type MemoryNode struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // fact, insight, pattern, tool_result, event, summary
	Content    string    `json:"content"`
	Importance float64   `json:"importance"`
	UsageCount int       `json:"usage_count"`
	Archived   bool      `json:"archived"`
	Metadata   string    `json:"metadata"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MemoryLink represents a relationship between memory nodes
type MemoryLink struct {
	ID           string    `json:"id"`
	SourceID     string    `json:"source_id"`
	TargetID     string    `json:"target_id"`
	RelationType string    `json:"relation_type"` // related_to, caused_by, solved_by, derived_from, part_of
	Metadata     string    `json:"metadata"`
	CreatedAt    time.Time `json:"created_at"`
}

// MemoryManager handles creation and linkage of knowledge artifacts
type MemoryManager struct {
	db *sql.DB
}

// NewMemoryManager creates a new MemoryManager
func NewMemoryManager(db *sql.DB) *MemoryManager {
	return &MemoryManager{
		db: db,
	}
}

// StoreFact creates a new fact node and returns its ID
func (mm *MemoryManager) StoreFact(ctx context.Context, content string, importance float64) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO memory_nodes (id, type, content, importance, metadata) VALUES (?, ?, ?, ?, ?)`
	_, err := mm.db.ExecContext(ctx, query, id, "fact", content, importance, "{}")
	return id, err
}

// StoreToolResult creates a node representing the output of a tool
func (mm *MemoryManager) StoreToolResult(ctx context.Context, toolName string, input string, result string) (string, error) {
	id := uuid.New().String()
	metadata := fmt.Sprintf(`{"tool": "%s", "input": "%s"}`, toolName, input)
	query := `INSERT INTO memory_nodes (id, type, content, importance, metadata) VALUES (?, ?, ?, ?, ?)`
	_, err := mm.db.ExecContext(ctx, query, id, "tool_result", result, 0.5, metadata)
	return id, err
}

// LinkMemories creates a relationship between two existing nodes
func (mm *MemoryManager) LinkMemories(ctx context.Context, sourceID string, targetID string, relation string) error {
	id := uuid.New().String()
	query := `INSERT INTO memory_links (id, source_id, target_id, relation_type) VALUES (?, ?, ?, ?)`
	_, err := mm.db.ExecContext(ctx, query, id, sourceID, targetID, relation)
	return err
}

// GetRelatedNodes retrieves nodes connected to the given node
func (mm *MemoryManager) GetRelatedNodes(ctx context.Context, nodeID string) ([]*MemoryNode, error) {
	query := `
		SELECT n.id, n.type, n.content, n.importance, n.usage_count, n.archived, n.metadata, n.created_at, n.updated_at
		FROM memory_nodes n
		JOIN memory_links l ON (l.source_id = n.id OR l.target_id = n.id)
		WHERE (l.source_id = ? OR l.target_id = ?) AND n.id != ?
	`
	rows, err := mm.db.QueryContext(ctx, query, nodeID, nodeID, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*MemoryNode
	for rows.Next() {
		var n MemoryNode
		err := rows.Scan(&n.ID, &n.Type, &n.Content, &n.Importance, &n.UsageCount, &n.Archived, &n.Metadata, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}
