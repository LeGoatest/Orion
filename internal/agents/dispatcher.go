package agents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"orion/internal/execution"
	"orion/internal/runtime"
)

type Dispatcher struct {
	Reg *Registry
	Sch *execution.Scheduler
	EB  *runtime.EventBus
}

func NewDispatcher(r *Registry, s *execution.Scheduler, eb *runtime.EventBus) *Dispatcher {
	return &Dispatcher{
		Reg: r,
		Sch: s,
		EB:  eb,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, stage, goalID string, task interface{}) error {
	agents := d.Reg.GetByCap(stage)
	if len(agents) == 0 {
		return fmt.Errorf("no agent for %s", stage)
	}

	d.Sch.Schedule(&AgentJob{
		IDVal:  uuid.New().String(),
		GoalID: goalID,
		Stage:  stage,
		Agent:  agents[0],
		Task:   task,
		EB:     d.EB,
	})
	return nil
}
