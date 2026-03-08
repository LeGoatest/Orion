package cognition

import (
	"context"
	"fmt"
	"time"

	"orion/internal/agents"
	"orion/internal/types"
)

type CognitionEngine struct {
	pipeline   OODALPipeline
	eventBus   *types.EventBus
	dispatcher *agents.Dispatcher
}

func NewCognitionEngine(p OODALPipeline, eb *types.EventBus, ad *agents.Dispatcher) *CognitionEngine {
	ce := &CognitionEngine{
		pipeline:   p,
		eventBus:   eb,
		dispatcher: ad,
	}

	// Subscribe to completion events to drive the multi-agent pipeline
	go ce.listenForEvents()

	return ce
}

func (ce *CognitionEngine) listenForEvents() {
	ch := ce.eventBus.Subscribe("*")
	for event := range ch {
		ce.HandleStageCompletion(context.Background(), event)
	}
}

func (ce *CognitionEngine) Process(ctx context.Context, goalID string, description string) error {
	fmt.Printf("CognitionEngine: Initializing goal %s\n", goalID)

	// First stage: OBSERVE
	return ce.dispatcher.Dispatch(ctx, "OBSERVE", goalID, description)
}

func (ce *CognitionEngine) HandleStageCompletion(ctx context.Context, event types.Event) {
	// Only handle stage completion events
	payload, ok := event.Payload.(map[string]interface{})
	if !ok { return }
	goalID, _ := payload["goal_id"].(string)

	switch event.Type {
	case "cognition.OBSERVE.completed":
		fmt.Printf("Pipeline [%s]: OBSERVE complete -> SYMBOL_LOOKUP\n", goalID)
		ce.dispatcher.Dispatch(ctx, "SYMBOL_LOOKUP", goalID, payload["result"])

	case "cognition.SYMBOL_LOOKUP.completed":
		fmt.Printf("Pipeline [%s]: SYMBOL_LOOKUP complete -> PATTERN_MATCH\n", goalID)
		ce.dispatcher.Dispatch(ctx, "PATTERN_MATCH", goalID, payload["result"])

	case "cognition.PATTERN_MATCH.completed":
		fmt.Printf("Pipeline [%s]: PATTERN_MATCH complete -> RETRIEVAL\n", goalID)
		ce.dispatcher.Dispatch(ctx, "RETRIEVAL", goalID, payload["result"])

	case "cognition.RETRIEVAL.completed":
		fmt.Printf("Pipeline [%s]: RETRIEVAL complete -> PLAN\n", goalID)
		ce.dispatcher.Dispatch(ctx, "PLAN", goalID, payload["result"])

	case "cognition.PLAN.completed":
		fmt.Printf("Pipeline [%s]: PLAN complete -> ACT\n", goalID)
		ce.dispatcher.Dispatch(ctx, "ACT", goalID, payload["result"])

	case "cognition.ACT.completed":
		fmt.Printf("Pipeline [%s]: ACT complete -> LEARN\n", goalID)
		ce.eventBus.Publish(types.Event{
			Type:      "cognition.LEARN.completed",
			Payload:   map[string]string{"goal_id": goalID},
			CreatedAt: time.Now(),
		})

	case "cognition.LEARN.completed":
		fmt.Printf("Pipeline [%s]: LEARN complete -> GARDEN (scheduled)\n", goalID)
		ce.eventBus.Publish(types.Event{
			Type:      "cognition.GARDEN.scheduled",
			Payload:   map[string]string{"goal_id": goalID},
			CreatedAt: time.Now(),
		})
		ce.eventBus.Publish(types.Event{
			Type:      "cognition.goal.completed",
			Payload:   map[string]string{"goal_id": goalID},
			CreatedAt: time.Now(),
		})
	}
}

type OODALPipeline interface {
	Observe(context.Context, interface{}) (interface{}, error)
	Orient(context.Context, interface{}) (interface{}, error)
	Decide(context.Context, interface{}) (interface{}, error)
	Act(context.Context, interface{}) (interface{}, error)
	Learn(context.Context, interface{}) error
	Garden(context.Context, string) error
}

type DefaultPipeline struct{}
func (p *DefaultPipeline) Observe(ctx context.Context, i interface{}) (interface{}, error) { return i, nil }
func (p *DefaultPipeline) Orient(ctx context.Context, i interface{}) (interface{}, error)  { return i, nil }
func (p *DefaultPipeline) Decide(ctx context.Context, i interface{}) (interface{}, error)  { return i, nil }
func (p *DefaultPipeline) Act(ctx context.Context, i interface{}) (interface{}, error)     { return i, nil }
func (p *DefaultPipeline) Learn(ctx context.Context, i interface{}) error                 { return nil }
func (p *DefaultPipeline) Garden(ctx context.Context, s string) error                     { return nil }
