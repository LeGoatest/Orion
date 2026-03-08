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
	CodeSymbols       []symbols.Symbol
	SymbolDefinitions []string
	CallGraphEdges    []string
	RelevantEvents    []string
	Nodes             []ScoredNode
}

type ScoredNode struct {
	ID        string
	Type      string
	Content   string
	Score     float64
	Metadata  map[string]interface{}
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

	// 1. Symbol lookup
	re.eb.Publish(types.Event{Type: "pipeline.symbol_lookup", CreatedAt: time.Now()})
	syms, _ := re.symbolQuery.ResolveSymbolReference(ctx, goal)

	// 2. Vector Search (Simulated)
	re.eb.Publish(types.Event{Type: "pipeline.vector_search", CreatedAt: time.Now()})
	vectorResults := re.simulateVectorSearch(goal)

	// 3. Graph Expansion
	re.eb.Publish(types.Event{Type: "pipeline.graph_expand", CreatedAt: time.Now()})
	graphResults := re.simulateGraphExpansion(vectorResults)

	// 4. Hybrid Scoring
	re.eb.Publish(types.Event{Type: "retrieval.ranked", CreatedAt: time.Now()})
	scoredNodes := re.scoreHybrid(vectorResults, graphResults)

	// Prepare bundle
	bundle := &ContextBundle{
		CodeSymbols: syms,
		Nodes:       scoredNodes,
	}

	re.eb.Publish(types.Event{Type: "pipeline.context_built", CreatedAt: time.Now()})
	return bundle, nil
}

func (re *RetrievalEngine) scoreHybrid(vector []ScoredNode, graph []ScoredNode) []ScoredNode {
	// Hybrid scoring formula:
	// score = 0.55 * semantic + 0.15 * graph + 0.15 * temporal + 0.10 * usage + 0.05 * importance

	weights := map[string]float64{
		"semantic":   0.55,
		"graph":      0.15,
		"temporal":   0.15,
		"usage":      0.10,
		"importance": 0.05,
	}

	scores := make(map[string]*ScoredNode)

	// Merge results
	for _, n := range vector {
		if _, ok := scores[n.ID]; !ok {
			scores[n.ID] = &n
		}
	}
	for _, n := range graph {
		if _, ok := scores[n.ID]; !ok {
			scores[n.ID] = &n
		}
	}

	finalNodes := make([]ScoredNode, 0, len(scores))
	for _, n := range scores {
		// Calculate components (simulated for now)
		semantic := n.Score // assume base score is semantic
		graphConn := re.calculateGraphConnectivity(n.ID)
		temporal := re.calculateTemporalRecency(n.Metadata["created_at"])
		usage := re.calculateUsageFrequency(n.ID)
		importance := re.calculateImportance(n.Metadata["priority"])

		n.Score = (weights["semantic"] * semantic) +
			(weights["graph"] * graphConn) +
			(weights["temporal"] * temporal) +
			(weights["usage"] * usage) +
			(weights["importance"] * importance)

		finalNodes = append(finalNodes, *n)
	}

	return finalNodes
}

func (re *RetrievalEngine) calculateGraphConnectivity(id string) float64 {
	// Simulated graph score [0-1]
	return 0.75
}

func (re *RetrievalEngine) calculateTemporalRecency(createdAt interface{}) float64 {
	// Simulated recency decay [0-1]
	return 0.9
}

func (re *RetrievalEngine) calculateUsageFrequency(id string) float64 {
	// Simulated usage score [0-1]
	return 0.4
}

func (re *RetrievalEngine) calculateImportance(priority interface{}) float64 {
	// Simulated importance score [0-1]
	return 0.5
}

func (re *RetrievalEngine) simulateVectorSearch(query string) []ScoredNode {
	return []ScoredNode{
		{ID: "node-1", Type: "memory", Content: "Related memory 1", Score: 0.85, Metadata: map[string]interface{}{"created_at": time.Now(), "priority": 1}},
		{ID: "node-2", Type: "code", Content: "Related code snippet", Score: 0.70, Metadata: map[string]interface{}{"created_at": time.Now().Add(-1 * time.Hour), "priority": 0.8}},
	}
}

func (re *RetrievalEngine) simulateGraphExpansion(seeds []ScoredNode) []ScoredNode {
	return []ScoredNode{
		{ID: "node-3", Type: "pattern", Content: "Linked pattern", Score: 0.60, Metadata: map[string]interface{}{"created_at": time.Now().Add(-24 * time.Hour), "priority": 0.5}},
	}
}
