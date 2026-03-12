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
	goal, err := Observe(ctx, e.client, intent)
	if err != nil {
		return err
	}

	c, err := Orient(ctx, e.client, goal)
	if err != nil {
		return err
	}

	plan, err := Decide(ctx, c)
	if err != nil {
		return err
	}

	err = Act(ctx, plan)
	if err != nil {
		return err
	}

	return Learn(ctx, e.client, goal)
}
