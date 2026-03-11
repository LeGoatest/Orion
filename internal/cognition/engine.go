package cognition

import (
	"context"
	"fmt"
	"orion/internal/agents"
	"orion/internal/storage/sqlite"
	"orion/internal/types"
)

type Engine struct {
	eb   *types.EventBus
	ad   *agents.Dispatcher
	pers *sqlite.Persistence
}

func NewEngine(eb *types.EventBus, ad *agents.Dispatcher) *Engine {
	ce := &Engine{eb: eb, ad: ad}
	go ce.listen()
	return ce
}

func (ce *Engine) SetPersistence(p *sqlite.Persistence) {
	ce.pers = p
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

func (ce *Engine) HandleCompletion(e types.Event) {
	payload, ok := e.Payload.(map[string]interface{})
	if !ok {
		return
	}
	gID, _ := payload["goal_id"].(string)

	dispatch := func(stage string, result interface{}) {
		if ce.pers != nil {
			ce.pers.TransitionStage(context.Background(), gID, stage, "System")
		}
		ce.ad.Dispatch(context.Background(), stage, gID, result)
	}

	switch e.Type {
	case "cognition.OBSERVE.completed":
		dispatch("SYMBOL_LOOKUP", payload["result"])
	case "cognition.SYMBOL_LOOKUP.completed":
		dispatch("PATTERN_MATCH", payload["result"])
	case "cognition.PATTERN_MATCH.completed":
		dispatch("RETRIEVAL", payload["result"])
	case "cognition.RETRIEVAL.completed":
		dispatch("PLAN", payload["result"])
	case "cognition.PLAN.completed":
		dispatch("ACT", payload["result"])
	case "cognition.ACT.completed":
		fmt.Printf("Goal %s completed successfully!\n", gID)
	}
}
