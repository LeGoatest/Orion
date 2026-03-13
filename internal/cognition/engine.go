package cognition

import (
	"context"
	"orion/ent"
)

type Engine struct {
	client *ent.Client
}

func NewEngine(client *ent.Client) *Engine {
	return &Engine{client: client}
}

func (e *Engine) Run(ctx context.Context, intent string) error {
	event, err := Observe(ctx, e.client, intent)
	if err != nil {
		return err
	}

	sm, err := Orient(ctx, e.client, event)
	if err != nil {
		return err
	}

	plan, err := Decide(ctx, sm)
	if err != nil {
		return err
	}

	err = Act(ctx, plan)
	if err != nil {
		return err
	}

	_, err = Learn(ctx, e.client, plan)
	return err
}
