package graph

import (
	"context"
)

type Ranking struct {
	expander *Expander
}

func NewRanking(e *Expander) *Ranking {
	return &Ranking{expander: e}
}

// Score calculates the graph relevance score for an expanded node
func (r *Ranking) Score(ctx context.Context, nodeID string, initialNodes []string) float64 {
	// A simple centrality-based score: nodes closer to initial search results get higher scores
	// For bootstrap, we return 1.0 for initial nodes and lower for expanded ones
	for _, id := range initialNodes {
		if id == nodeID {
			return 1.0
		}
	}

	// Check if it's a 1-hop neighbor
	neighbors, _ := r.expander.Expand(ctx, initialNodes, 1)
	for _, id := range neighbors {
		if id == nodeID {
			return 0.5
		}
	}

	return 0.25
}
