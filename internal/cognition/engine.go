package cognition

import (
	"context"
	"fmt"
	"time"

	"orion/internal/agents"
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
	dispatcher    *agents.Dispatcher
}

func NewCognitionEngine(pipeline OODALPipeline, eb *types.EventBus, re *retrieval.RetrievalEngine, pe *pattern.Engine, ad *agents.Dispatcher) *CognitionEngine {
	return &CognitionEngine{
		pipeline:      pipeline,
		eventBus:      eb,
		retrieval:     re,
		patternEngine: pe,
		dispatcher:    ad,
	}
}

func (ce *CognitionEngine) Process(ctx context.Context, g *goal.Goal) error {
	fmt.Printf("CognitionEngine: Starting Multi-Agent OODA-L loop for goal %s\n", g.ID)

	// 1. OBSERVE
	obs, err := ce.pipeline.Observe(ctx, g.Description)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "observation_recorded", Payload: obs, CreatedAt: time.Now()})

	// 2. ORIENT
	if p, found := ce.patternEngine.Match(ctx, g.Description); found {
		fmt.Printf("Cognition: Pattern match found. Bypassing planning.\n")
		return ce.patternEngine.ExecutePattern(ctx, p)
	}

	bundle, err := ce.retrieval.AssembleContext(ctx, g.Description)
	if err != nil { return err }

	// 3. DECIDE
	decision, err := ce.pipeline.Decide(ctx, bundle)
	if err != nil { return err }

	// 4. DISPATCH AGENTS (New Stage)
	fmt.Println("Cognition: Dispatching agents for ACT stage.")
	if err := ce.dispatcher.Dispatch(ctx, "general_task", decision); err != nil {
		fmt.Printf("Cognition: Dispatch failed: %v. Falling back to ACT.\n", err)
	}

	// 5. ACT (Default or Agent-led)
	result, err := ce.pipeline.Act(ctx, decision)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "tool_executed", Payload: result, CreatedAt: time.Now()})

	// 6. LEARN
	if err := ce.pipeline.Learn(ctx, result); err != nil { return err }

	// 7. GARDEN
	if err := ce.pipeline.Garden(ctx, g.ID); err != nil { return err }

	return nil
}
