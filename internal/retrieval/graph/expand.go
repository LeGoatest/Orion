package graph

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Expander struct {
	Db *sql.DB
}

func NewExpander(db *sql.DB) *Expander {
	return &Expander{Db: db}
}

// Expand retrieves 1-hop and 2-hop neighbors for a given set of node IDs
func (e *Expander) Expand(ctx context.Context, nodeIDs []string, depth int) ([]string, error) {
	if len(nodeIDs) == 0 || depth <= 0 {
		return nodeIDs, nil
	}

	currentSet := make(map[string]bool)
	for _, id := range nodeIDs {
		currentSet[id] = true
	}

	for d := 0; d < depth; d++ {
		// Build placeholders for the query
		placeholders := make([]string, len(nodeIDs))
		args := make([]interface{}, len(nodeIDs))
		for i, id := range nodeIDs {
			placeholders[i] = "?"
			args[i] = id
		}

		query := fmt.Sprintf(`
			SELECT target_id FROM memory_links WHERE source_id IN (%s)
			UNION
			SELECT source_id FROM memory_links WHERE target_id IN (%s)
		`, strings.Join(placeholders, ","), strings.Join(placeholders, ","))

		// Duplicate args for the second half of UNION
		fullArgs := append(args, args...)

		rows, err := e.Db.QueryContext(ctx, query, fullArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to query graph neighbors: %w", err)
		}
		defer rows.Close()

		var nextNodeIDs []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err == nil {
				if !currentSet[id] {
					currentSet[id] = true
					nextNodeIDs = append(nextNodeIDs, id)
				}
			}
		}

		if len(nextNodeIDs) == 0 {
			break
		}
		nodeIDs = nextNodeIDs
	}

	result := make([]string, 0, len(currentSet))
	for id := range currentSet {
		result = append(result, id)
	}

	return result, nil
}
