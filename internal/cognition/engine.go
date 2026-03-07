package cognition

import (
	"context"
	"fmt"
	"time"

	"orion/internal/retrieval"
	"orion/internal/runtime/goal"
	"orion/internal/types"
)

type CognitionEngine struct {
	pipeline     OODALPipeline
	eventBus     *types.EventBus
	retrieval    *retrieval.RetrievalEngine
}

func NewCognitionEngine(pipeline OODALPipeline, eb *types.EventBus, re *retrieval.RetrievalEngine) *CognitionEngine {
	return &CognitionEngine{
		pipeline:  pipeline,
		eventBus:  eb,
		retrieval: re,
	}
}

func (ce *CognitionEngine) Process(ctx context.Context, g *goal.Goal) error {
	fmt.Printf("CognitionEngine: Starting OODA-L-G loop for goal %s\n", g.ID)

	// OBSERVE
	obs, err := ce.pipeline.Observe(ctx, g.Description)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "observation_recorded", Payload: obs, CreatedAt: time.Now()})

	// ORIENT
	bundle, err := ce.pipeline.Orient(ctx, obs)
	if err != nil { return err }

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

func (p *DefaultPipeline) Observe(ctx context.Context, input interface{}) (interface{}, error) {
	fmt.Println("Cognition: Phase Observe - Creating memory node for intent")
	return input, nil
}

func (p *DefaultPipeline) Orient(ctx context.Context, observation interface{}) (interface{}, error) {
	fmt.Println("Cognition: Phase Orient - Assembling ContextBundle")
	return &retrieval.ContextBundle{}, nil
}

func (p *DefaultPipeline) Decide(ctx context.Context, orientation interface{}) (interface{}, error) {
	fmt.Println("Cognition: Phase Decide - Generating execution plan")
	return "plan", nil
}

func (p *DefaultPipeline) Act(ctx context.Context, decision interface{}) (interface{}, error) {
	fmt.Println("Cognition: Phase Act - Executing tools")
	return "result", nil
}

func (p *DefaultPipeline) Learn(ctx context.Context, result interface{}) error {
	fmt.Println("Cognition: Phase Learn - Updating knowledge graph")
	return nil
}

func (p *DefaultPipeline) Garden(ctx context.Context, goalID string) error {
	fmt.Println("Cognition: Phase Garden - Maintaining memory quality")
	return nil
}
