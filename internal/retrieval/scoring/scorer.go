package scoring

import (
	"orion/internal/storage/sqlite"
)

type Scorer struct {
	db *sqlite.DB
}

func NewScorer(db *sqlite.DB) *Scorer {
	return &Scorer{db: db}
}

type Factors struct {
	Semantic   float32
	Graph      float32
	Temporal   float32
	Usage      float32
	Importance float32
}

func (s *Scorer) Calculate(f Factors) float32 {
	// Orion Hybrid Scoring formula:
	// score = 0.55 semantic + 0.15 graph + 0.15 temporal + 0.10 usage + 0.05 importance
	return (0.55 * f.Semantic) +
		(0.15 * f.Graph) +
		(0.15 * f.Temporal) +
		(0.10 * f.Usage) +
		(0.05 * f.Importance)
}
