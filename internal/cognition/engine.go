package cognition

import (
	"context"
	"fmt"
	"orion/ent"
)

type Engine struct {
	client *ent.Client
	gov    *GovernanceValidator
}

func NewEngine(client *ent.Client) *Engine {
	return &Engine{
		client: client,
		gov:    &GovernanceValidator{},
	}
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

	// SAGE Governance Gate
	vp, err := e.gov.Validate(ctx, plan, sm)
	if err != nil {
		return err
	}

	fmt.Printf("Governance: Plan disposition is %s. Reason: %s\n", vp.Disposition, vp.ValidationReason)

	if vp.Disposition != "approved" {
		fmt.Printf("Governance: Blocking execution of goal %s\n", plan.GoalID)
		return fmt.Errorf("governance refusal: %s", vp.ValidationReason)
	}

	err = Act(ctx, vp)
	if err != nil {
		return err
	}

	_, err = Learn(ctx, e.client, plan)
	return err
}
