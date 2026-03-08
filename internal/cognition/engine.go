package cognition

import (
	"context"
	"fmt"
	"time"

	"orion/internal/pattern"
	"orion/internal/retrieval"
	"orion/internal/runtime/goal"
	"orion/internal/types"
)

type CognitionEngine struct {
	pipeline      OODALPipeline
	eventBus      *types.EventBus
	retrieval     *retrieval.RetrievalEngine
	patternEngine *pattern.Engine
}

func NewCognitionEngine(pipeline OODALPipeline, eb *types.EventBus, re *retrieval.RetrievalEngine, pe *pattern.Engine) *CognitionEngine {
	return &CognitionEngine{
		pipeline:      pipeline,
		eventBus:      eb,
		retrieval:     re,
		patternEngine: pe,
	}
}

func (ce *CognitionEngine) Process(ctx context.Context, g *goal.Goal) error {
	fmt.Printf("CognitionEngine: Starting Optimized OODA-L loop for goal %s\n", g.ID)

	// 1. OBSERVE
	obs, err := ce.pipeline.Observe(ctx, g.Description)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "observation_recorded", Payload: obs, CreatedAt: time.Now()})

	// 2. ORIENT (Optimized Pipeline)
	// Phase 2a: Symbol Lookup & Pattern Matching occur inside RetrievalEngine.AssembleContext
	// or are orchestrated here. According to specs, Orient stage flow is:
	// symbol lookup -> pattern match -> embedding -> vector retrieval -> graph expansion -> hybrid scoring

	// We perform Pattern Match here for potential DECIDE bypass.
	ce.eventBus.Publish(types.Event{Type: "pipeline.pattern_match", CreatedAt: time.Now()})
	if p, found := ce.patternEngine.Match(ctx, g.Description); found {
		fmt.Printf("Cognition: Strong pattern match found! Bypassing complex planning for pattern: %s\n", p.ID)
		ce.eventBus.Publish(types.Event{Type: "pipeline.pattern_match_success", Payload: p.ID, CreatedAt: time.Now()})
		return ce.patternEngine.ExecutePattern(ctx, p)
	}

	bundle, err := ce.retrieval.AssembleContext(ctx, g.Description)
	if err != nil { return err }

	// 3. DECIDE
	decision, err := ce.pipeline.Decide(ctx, bundle)
	if err != nil { return err }

	// 4. ACT
	result, err := ce.pipeline.Act(ctx, decision)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "tool_executed", Payload: result, CreatedAt: time.Now()})

	// 5. LEARN
	if err := ce.pipeline.Learn(ctx, result); err != nil { return err }

	// 6. GARDEN
	if err := ce.pipeline.Garden(ctx, g.ID); err != nil { return err }

	return nil
}
