package agents

import ("context"; "fmt"; "orion/internal/execution"; "orion/internal/types"; "time"; "github.com/google/uuid")

type Dispatcher struct { Reg *Registry; Sch *execution.Scheduler; EB *types.EventBus }
func NewDispatcher(r *Registry, s *execution.Scheduler, eb *types.EventBus) *Dispatcher { return &Dispatcher{Reg: r, Sch: s, EB: eb} }
func (d *Dispatcher) Dispatch(ctx context.Context, stage, goalID string, task interface{}) error { m := d.Reg.GetByCap(stage); if len(m) == 0 { return fmt.Errorf("no agent for %s", stage) }; d.Sch.Schedule(&AgentJob{IDVal: uuid.New().String(), GoalID: goalID, Stage: stage, Agent: m[0], Task: task, EB: d.EB}); return nil }
