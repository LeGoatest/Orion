package cognition

import (
	"context"
	"fmt"
	"orion/internal/agents"
	"orion/internal/runtime"
)

type Engine struct {
	eb *runtime.EventBus
	ad *agents.Dispatcher
}

func NewEngine(eb *runtime.EventBus, ad *agents.Dispatcher) *Engine {
	ce := &Engine{eb: eb, ad: ad}
	go ce.listen()
	return ce
}

func (ce *Engine) listen() {
	ch := ce.eb.Subscribe("*")
	for e := range ch {
		ce.HandleCompletion(e)
	}
}

func (ce *Engine) Process(ctx context.Context, id, desc string) error {
	fmt.Printf("Pipeline starting for goal %s\n", id)
	return ce.ad.Dispatch(ctx, "OBSERVE", id, desc)
}

func (ce *Engine) HandleCompletion(e runtime.Event) {
	payload, ok := e.Payload.(map[string]interface{})
	if !ok {
		return
	}
	gID, _ := payload["goal_id"].(string)
	switch e.Type {
	case "cognition.OBSERVE.completed":
		ce.ad.Dispatch(context.Background(), "SYMBOL_LOOKUP", gID, payload["result"])
	case "cognition.SYMBOL_LOOKUP.completed":
		ce.ad.Dispatch(context.Background(), "PATTERN_MATCH", gID, payload["result"])
	case "cognition.PATTERN_MATCH.completed":
		ce.ad.Dispatch(context.Background(), "RETRIEVAL", gID, payload["result"])
	case "cognition.RETRIEVAL.completed":
		ce.ad.Dispatch(context.Background(), "PLAN", gID, payload["result"])
	case "cognition.PLAN.completed":
		ce.ad.Dispatch(context.Background(), "ACT", gID, payload["result"])
	case "cognition.ACT.completed":
		fmt.Printf("Goal %s completed successfully!\n", gID)
	}
}
