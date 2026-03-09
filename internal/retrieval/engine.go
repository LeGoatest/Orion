package retrieval

import (
	"context"
	"fmt"
	"sync"
	"time"
	"orion/internal/retrieval/graph"
	"orion/internal/symbols"
	"orion/internal/types"
	"orion/internal/retrieval/scoring"
)

type ContextBundle struct {
	Facts             []string
	Patterns          []string
	Insights          []string
	CodeSymbols       []symbols.Symbol
	CallGraphEdges    []string
	RelevantEvents    []string
	FormattedContext  string
}

type RetrievalEngine struct {
	eb          *types.EventBus
	symbolQuery *symbols.Query
	graphExpand *graph.Expander
	scorer      *scoring.Scorer
}

func NewRetrievalEngine(eb *types.EventBus, sq *symbols.Query, ge *graph.Expander) *RetrievalEngine {
	return &RetrievalEngine{
		eb:          eb,
		symbolQuery: sq,
		graphExpand: ge,
		scorer:      &scoring.Scorer{},
	}
}

func (re *RetrievalEngine) AssembleContext(ctx context.Context, goal string) (*ContextBundle, error) {
	fmt.Printf("Retrieval Engine: Assembling Coordinated Context for: %s\n", goal)

	// 1. Symbol lookup
	re.eb.Publish(types.Event{Type: "pipeline.symbol_lookup", CreatedAt: time.Now()})
	syms, _ := re.symbolQuery.ResolveSymbolReference(ctx, goal)

	// 2. Concurrent Vector search and Graph Expansion
	var candidates []scoring.RetrievalCandidate
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		re.eb.Publish(types.Event{Type: "pipeline.vector_search", CreatedAt: time.Now()})
		// Simulated vector results
		mu.Lock()
		candidates = append(candidates, scoring.RetrievalCandidate{ID: "mem-1", Semantic: 0.9, Content: "Relevant fact from vector search"})
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		re.eb.Publish(types.Event{Type: "pipeline.graph_expand", CreatedAt: time.Now()})
		ids, _ := re.graphExpand.Expand(ctx, []string{"root-node"}, 2)
		mu.Lock()
		for _, id := range ids {
			candidates = append(candidates, scoring.RetrievalCandidate{ID: id, Graph: 0.8, Content: "Graph neighbor node content"})
		}
		mu.Unlock()
	}()

	wg.Wait()

	// 3. Hybrid scoring
	re.eb.Publish(types.Event{Type: "retrieval.ranked", CreatedAt: time.Now()})
	ranked := re.scorer.RankCandidates(ctx, candidates)

	// 4. Build deterministic bundle
	bundle := &ContextBundle{
		CodeSymbols: syms,
		Facts:       []string{},
	}
	for _, c := range ranked {
		bundle.Facts = append(bundle.Facts, c.Content)
	}

	re.eb.Publish(types.Event{Type: "pipeline.context_built", CreatedAt: time.Now()})
	return bundle, nil
}
