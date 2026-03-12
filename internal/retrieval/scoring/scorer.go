package retrieval

import (
	"context"
)

type RetrievalCandidate struct {
	ID         string
	Content    string
	Semantic   float32
	Graph      float32
	Temporal   float32
	Usage      float32
	Importance float32
}

type Scorer struct{}

func (s *Scorer) RankCandidates(ctx context.Context, candidates []RetrievalCandidate) []RetrievalCandidate {
	// Simple pass-through or basic sort for now to fix build
	return candidates
}
