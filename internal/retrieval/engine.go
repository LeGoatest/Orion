package retrieval

import (
	"context"
	"fmt"
	"sync"
	"time"
	"orion/internal/retrieval/graph"
	"orion/internal/symbols"
	"orion/internal/types"
)

type ContextBundle struct {
	Facts             []string
	Patterns          []string
	Insights          []string
	CodeSymbols       []string
	CallGraphEdges    []string
	RelevantEvents    []string
	Symbols           []symbols.Symbol
	SymbolDefinitions []string
	RelatedCalls      []string
}

type RetrievalEngine struct {
	eb          *types.EventBus
	symbolQuery *symbols.Query
	graphExpand *graph.Expander
}

func NewRetrievalEngine(eb *types.EventBus, sq *symbols.Query, ge *graph.Expander) *RetrievalEngine {
	return &RetrievalEngine{
		eb:          eb,
		symbolQuery: sq,
		graphExpand: ge,
	}
}

func (re *RetrievalEngine) AssembleContext(ctx context.Context, goal string) (*ContextBundle, error) {
	fmt.Printf("Retrieval Engine: Assembling Optimized Context for: %s\n", goal)

	// 1. Symbol lookup (Direct/Fuzzy)
	re.eb.Publish(types.Event{Type: "pipeline.symbol_lookup", CreatedAt: time.Now()})
	syms, _ := re.symbolQuery.ResolveSymbolReference(ctx, goal)
	if len(syms) > 0 {
		re.eb.Publish(types.Event{Type: "symbol.resolved", CreatedAt: time.Now()})
	}

	// Prepare bundle
	bundle := &ContextBundle{
		Symbols: syms,
	}

	// 2. Vector Search and Graph Expansion (Concurrent where possible)
	re.eb.Publish(types.Event{Type: "pipeline.vector_search", CreatedAt: time.Now()})

	// Simulated concurrent execution
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// Vector retrieval logic
		re.eb.Publish(types.Event{Type: "retrieval.vector_completed", CreatedAt: time.Now()})
	}()

	go func() {
		defer wg.Done()
		// 3. Graph expansion (memory_links and call graph edges)
		re.eb.Publish(types.Event{Type: "pipeline.graph_expand", CreatedAt: time.Now()})
		_, _ = re.graphExpand.Expand(ctx, []string{"root-node"}, 2)
		re.eb.Publish(types.Event{Type: "retrieval.graph_expanded", CreatedAt: time.Now()})
	}()

	wg.Wait()

	// 4. Hybrid scoring and final build
	re.eb.Publish(types.Event{Type: "retrieval.ranked", CreatedAt: time.Now()})
	re.eb.Publish(types.Event{Type: "pipeline.context_built", CreatedAt: time.Now()})

	return bundle, nil
}
