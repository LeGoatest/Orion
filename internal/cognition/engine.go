package cognition

import (
	"context"
	"fmt"
	"orion/internal/runtime"
	"sync"
)

type Event types.Event // Type alias for convenience if needed, or import types directly

type Engine struct {
	eb  *runtime.EventBus
	ctx context.Context
	wg  sync.WaitGroup
}

func NewEngine(eb *runtime.EventBus, ctx context.Context) *Engine {
	return &Engine{eb: eb, ctx: ctx}
}

func (ce *Engine) Start() {
	fmt.Println("Cognition Engine: Starting...")
}
