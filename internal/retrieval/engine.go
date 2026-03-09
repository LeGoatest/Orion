package retrieval

import (
	"context"
	"fmt"
	"orion/internal/retrieval/scoring"
	"orion/internal/runtime"
	"orion/internal/symbols"
	"sync"
	"time"
)

type ContextBundle struct {
	Facts            []string
	Patterns         []string
	Insights         []string
	CodeSymbols      []symbols.Symbol
	CallGraphEdges   []string
	RelevantEvents   []string
	FormattedContext string
}

type RetrievalEngine struct {
	eb          *runtime.EventBus
	symbolQuery *symbols.Query
	scorer      *scoring.Scorer
}

func NewRetrievalEngine(eb *runtime.EventBus, sq *symbols.Query) *RetrievalEngine {
	return &RetrievalEngine{
		eb:          eb,
		symbolQuery: sq,
		scorer:      &scoring.Scorer{},
	}
}

func (re *RetrievalEngine) AssembleContext(ctx context.Context, goal string) (*ContextBundle, error) {
	fmt.Printf("Retrieval Engine: Assembling Coordinated Context for: %s\n", goal)

	// 1. Symbol lookup
	re.eb.Publish(runtime.Event{Type: "pipeline.symbol_lookup", CreatedAt: time.Now()})
	syms, _ := re.symbolQuery.ResolveSymbolReference(ctx, goal)

	// 2. Concurrent Vector search and Graph Expansion
	var candidates []scoring.RetrievalCandidate
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		re.eb.Publish(runtime.Event{Type: "pipeline.vector_search", CreatedAt: time.Now()})
		// Simulated vector results
		mu.Lock()
		candidates = append(candidates, scoring.RetrievalCandidate{ID: "mem-1", Semantic: 0.9, Content: "Relevant fact from vector search"})
		mu.Unlock()
	}()

	wg.Wait()

	// 3. Hybrid scoring
	re.eb.Publish(runtime.Event{Type: "retrieval.ranked", CreatedAt: time.Now()})
	ranked := re.scorer.RankCandidates(ctx, candidates)

	// 4. Build deterministic bundle
	bundle := &ContextBundle{
		CodeSymbols: syms,
		Facts:       []string{},
	}
	for _, c := range ranked {
		bundle.Facts = append(bundle.Facts, c.Content)
	}

	re.eb.Publish(runtime.Event{Type: "pipeline.context_built", CreatedAt: time.Now()})
	return bundle, nil
}
