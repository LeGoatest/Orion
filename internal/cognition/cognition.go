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

	// Subscribe to completion events to drive the pipeline
	go ce.listenForEvents()

	return ce
}

func (ce *CognitionEngine) listenForEvents() {
	ch := ce.eventBus.Subscribe("*")
	for event := range ch {
		ce.HandleStageCompletion(context.Background(), event)
	}
}

func (ce *CognitionEngine) Process(ctx context.Context, goal interface{}) error {
	fmt.Println("CognitionEngine: Starting coordination chain.")
	ce.eventBus.Publish(types.Event{Type: "cognition.observe.completed", Payload: goal, CreatedAt: time.Now()})
	return nil
}

func (ce *CognitionEngine) HandleStageCompletion(ctx context.Context, event types.Event) {
	switch event.Type {
	case "cognition.observe.completed":
		fmt.Println("Pipeline: Observe complete -> SymbolLookup")
		ce.dispatcher.Dispatch(ctx, "symbol_lookup", event.Payload)
	case "cognition.symbol_lookup.completed":
		fmt.Println("Pipeline: SymbolLookup complete -> PatternMatch")
		ce.dispatcher.Dispatch(ctx, "pattern_match", event.Payload)
	case "cognition.pattern.completed":
		fmt.Println("Pipeline: PatternMatch complete -> Retrieval")
		ce.dispatcher.Dispatch(ctx, "vector_retrieval", event.Payload)
	case "cognition.retrieval.completed":
		fmt.Println("Pipeline: Retrieval complete -> Planning")
		ce.dispatcher.Dispatch(ctx, "execution_planning", event.Payload)
	case "cognition.plan.completed":
		fmt.Println("Pipeline: Planning complete -> Act")
		ce.dispatcher.Dispatch(ctx, "act_stage", event.Payload)
	case "cognition.act.completed":
		fmt.Println("Pipeline: Act complete -> Learn")
		ce.eventBus.Publish(types.Event{Type: "cognition.learn.completed", CreatedAt: time.Now()})
	}
}

type OODALPipeline interface { Observe(context.Context, interface{}) (interface{}, error); Orient(context.Context, interface{}) (interface{}, error); Decide(context.Context, interface{}) (interface{}, error); Act(context.Context, interface{}) (interface{}, error); Learn(context.Context, interface{}) error; Garden(context.Context, string) error }

type DefaultPipeline struct{}
func (p *DefaultPipeline) Observe(ctx context.Context, i interface{}) (interface{}, error) { return i, nil }
func (p *DefaultPipeline) Orient(ctx context.Context, i interface{}) (interface{}, error)  { return i, nil }
func (p *DefaultPipeline) Decide(ctx context.Context, i interface{}) (interface{}, error)  { return i, nil }
func (p *DefaultPipeline) Act(ctx context.Context, i interface{}) (interface{}, error)     { return i, nil }
func (p *DefaultPipeline) Learn(ctx context.Context, i interface{}) error                 { return nil }
func (p *DefaultPipeline) Garden(ctx context.Context, s string) error                     { return nil }
