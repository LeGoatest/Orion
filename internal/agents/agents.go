package agents

import ("context"; "fmt"; "orion/internal/types"; "time"; "orion/internal/execution"; "github.com/google/uuid"; "sync")

type AutonomousAgent interface { Name() string; Capabilities() []string; ExecuteTask(context.Context, interface{}) (interface{}, error) }
type AgentJob struct { IDVal, GoalID, Stage string; Agent AutonomousAgent; Task interface{}; EB *types.EventBus }
func (j *AgentJob) ID() string { return j.IDVal }
func (j *AgentJob) Type() string { return "AgentTask" }
func (j *AgentJob) Execute(ctx context.Context) error {
	j.EB.Publish(types.Event{Type: "agent.task_received", GoalID: j.GoalID, Stage: j.Stage, Payload: j.Agent.Name(), CreatedAt: time.Now()})
	r, err := j.Agent.ExecuteTask(ctx, j.Task)
	if err != nil { j.EB.Publish(types.Event{Type: "agent.error", GoalID: j.GoalID, Stage: j.Stage, Payload: map[string]string{"agent": j.Agent.Name(), "error": err.Error()}, CreatedAt: time.Now()}); return err }
	j.EB.Publish(types.Event{Type: fmt.Sprintf("cognition.%s.completed", j.Stage), GoalID: j.GoalID, Stage: j.Stage, Payload: map[string]interface{}{"goal_id": j.GoalID, "result": r}, CreatedAt: time.Now()})
	return nil
}
type Registry struct { mu sync.RWMutex; m map[string]AutonomousAgent }
func NewRegistry() *Registry { return &Registry{m: make(map[string]AutonomousAgent)} }
func (r *Registry) RegisterAgent(a AutonomousAgent) { r.mu.Lock(); defer r.mu.Unlock(); r.m[a.Name()] = a }
func (r *Registry) GetByCap(c string) []AutonomousAgent { r.mu.RLock(); defer r.mu.RUnlock(); var res []AutonomousAgent; for _, a := range r.m { for _, cap := range a.Capabilities() { if cap == c { res = append(res, a); break } } }; return res }
func (r *Registry) List() []AutonomousAgent { r.mu.RLock(); defer r.mu.RUnlock(); var res []AutonomousAgent; for _, a := range r.m { res = append(res, a) }; return res }
type Supervisor struct { Reg *Registry; EB *types.EventBus }
func NewSupervisor(r *Registry, eb *types.EventBus) *Supervisor { return &Supervisor{Reg: r, EB: eb} }
func (s *Supervisor) StartAgents(ctx context.Context) error { for _, a := range s.Reg.List() { s.EB.Publish(types.Event{Type: "agent.started", Payload: a.Name(), CreatedAt: time.Now()}) }; return nil }
type Dispatcher struct { Reg *Registry; Sch *execution.Scheduler; EB *types.EventBus }
func NewDispatcher(r *Registry, s *execution.Scheduler, eb *types.EventBus) *Dispatcher { return &Dispatcher{Reg: r, Sch: s, EB: eb} }
func (d *Dispatcher) Dispatch(ctx context.Context, stage, goalID string, task interface{}) error { m := d.Reg.GetByCap(stage); if len(m) == 0 { return fmt.Errorf("no agent for %s", stage) }; d.Sch.Schedule(&AgentJob{IDVal: uuid.New().String(), GoalID: goalID, Stage: stage, Agent: m[0], Task: task, EB: d.EB}); return nil }
