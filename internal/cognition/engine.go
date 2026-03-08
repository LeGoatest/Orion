package cognition

import (
	"context"
	"fmt"
	"time"

	"orion/internal/agents"
	"orion/internal/runtime/goal"
	"orion/internal/types"
)

type CognitionEngine struct {
	pipeline   OODALPipeline
	eventBus   *types.EventBus
	dispatcher *agents.Dispatcher
}

func NewCognitionEngine(pipeline OODALPipeline, eb *types.EventBus, ad *agents.Dispatcher) *CognitionEngine {
	return &CognitionEngine{
		pipeline:   pipeline,
		eventBus:   eb,
		dispatcher: ad,
	}
}

// Process initiates the coordinated multi-agent cognition flow
func (ce *CognitionEngine) Process(ctx context.Context, g *goal.Goal) error {
	fmt.Printf("CognitionEngine: Starting Coordinated Pipeline for goal %s\n", g.ID)

	// OBSERVE
	obs, err := ce.pipeline.Observe(ctx, g.Description)
	if err != nil { return err }
	ce.eventBus.Publish(types.Event{Type: "cognition.observe.completed", Payload: obs, CreatedAt: time.Now()})

	// Dispatch next stages via events/dispatcher
	// The flow is now event-driven:
	// observe.completed -> symbol_lookup
	// symbol_lookup.completed -> pattern_match
	// ... and so on.

	// Start the chain
	return ce.dispatcher.Dispatch(ctx, "symbol_resolution", g.ID)
}

func (ce *CognitionEngine) HandleStageCompletion(ctx context.Context, event types.Event) {
	// Logic to trigger next stage based on event.Type
	switch event.Type {
	case "cognition.observe.completed":
		ce.dispatcher.Dispatch(ctx, "symbol_resolution", event.Payload)
	case "cognition.symbol_lookup.completed":
		ce.dispatcher.Dispatch(ctx, "pattern_matching", event.Payload)
	case "cognition.pattern_match.completed":
		ce.dispatcher.Dispatch(ctx, "vector_search", event.Payload)
	case "cognition.retrieval.completed":
		ce.dispatcher.Dispatch(ctx, "planning", event.Payload)
	case "cognition.plan.completed":
		ce.dispatcher.Dispatch(ctx, "act_stage_execution", event.Payload)
	case "cognition.act.completed":
		ce.eventBus.Publish(types.Event{Type: "cognition.learn.completed", CreatedAt: time.Now()})
	}
}
