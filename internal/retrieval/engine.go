package retrieval

type ContextBundle struct {
	Facts            []string
	Patterns         []string
	Insights         []string
	CodeSymbols      []string
	CallGraphEdges   []string
	RelevantEvents   []string
}

type RetrievalEngine struct{}

func NewRetrievalEngine() *RetrievalEngine {
	return &RetrievalEngine{}
}

func (re *RetrievalEngine) AssembleContext(goal string) (*ContextBundle, error) {
	return &ContextBundle{}, nil
}
