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
	fmt.Printf("CognitionEngine: Starting OODA-L-G loop for goal %s\n", g.ID)

	// OBSERVE
	obs, err := ce.pipeline.Observe(ctx, g.Description)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "observation_recorded", Payload: obs, CreatedAt: time.Now()})

	// ORIENT
	bundle, err := ce.retrieval.AssembleContext(ctx, g.Description)
	if err != nil { return err }

	// PATTERN MATCH - Bypasses DECIDE if strong match found
	if p, found := ce.patternEngine.Match(ctx, g.Description); found {
		fmt.Printf("Cognition: Strong pattern match found! Bypassing complex planning for pattern: %s\n", p.ID)
		return ce.patternEngine.ExecutePattern(ctx, p)
	}

	// DECIDE
	decision, err := ce.pipeline.Decide(ctx, bundle)
	if err != nil { return err }

	// ACT
	result, err := ce.pipeline.Act(ctx, decision)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "tool_executed", Payload: result, CreatedAt: time.Now()})

	// LEARN
	if err := ce.pipeline.Learn(ctx, result); err != nil { return err }

	// GARDEN
	if err := ce.pipeline.Garden(ctx, g.ID); err != nil { return err }

	return nil
}
