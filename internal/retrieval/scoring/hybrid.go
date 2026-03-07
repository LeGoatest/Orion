package scoring

type RetrievalCandidate struct {
	ID         string
	Semantic   float64
	Graph      float64
	Temporal   float64
	Usage      float64
	Importance float64
}

func CalculateHybridScore(c RetrievalCandidate) float64 {
	// score = 0.55 * semantic + 0.15 * graph + 0.15 * temporal + 0.10 * usage + 0.05 * importance
	return (0.55 * c.Semantic) +
		(0.15 * c.Graph) +
		(0.15 * c.Temporal) +
		(0.10 * c.Usage) +
		(0.05 * c.Importance)
}
