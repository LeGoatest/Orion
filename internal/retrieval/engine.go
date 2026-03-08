package retrieval

import (
	"context"
	"fmt"
	"time"
	"orion/internal/retrieval/graph"
	"orion/internal/symbols"
	"orion/internal/types"
)

type ContextBundle struct {
	Facts            []string
	Patterns         []string
	Insights         []string
	CodeSymbols      []string
	CallGraphEdges   []string
	RelevantEvents   []string
	Symbols          []symbols.Symbol
	SymbolDefinitions []string
	RelatedCalls     []string
}

type RetrievalEngine struct {
	eb           *types.EventBus
	symbolQuery  *symbols.Query
	graphExpand  *graph.Expander
}

func NewRetrievalEngine(eb *types.EventBus, sq *symbols.Query, ge *graph.Expander) *RetrievalEngine {
	return &RetrievalEngine{
		eb:          eb,
		symbolQuery: sq,
		graphExpand: ge,
	}
}

func (re *RetrievalEngine) AssembleContext(ctx context.Context, goal string) (*ContextBundle, error) {
	fmt.Printf("Retrieval Engine: Assembling context for: %s\n", goal)

	// 1. Symbol lookup first
	re.eb.Publish(types.Event{Type: "symbol.lookup", CreatedAt: time.Now()})
	syms, _ := re.symbolQuery.ResolveSymbolReference(ctx, goal)
	if len(syms) > 0 {
		re.eb.Publish(types.Event{Type: "symbol.resolved", CreatedAt: time.Now()})
	} else {
		re.eb.Publish(types.Event{Type: "symbol.miss", CreatedAt: time.Now()})
	}

	// 2. Vector search (placeholder logic)
	re.eb.Publish(types.Event{Type: "retrieval.vector_completed", CreatedAt: time.Now()})

	// 3. Graph expansion
	_, _ = re.graphExpand.Expand(ctx, []string{"root-node"}, 2)
	re.eb.Publish(types.Event{Type: "retrieval.graph_expanded", CreatedAt: time.Now()})

	// 4. Hybrid scoring and bundle build
	re.eb.Publish(types.Event{Type: "retrieval.ranked", CreatedAt: time.Now()})

	return &ContextBundle{
		Symbols: syms,
	}, nil
}
