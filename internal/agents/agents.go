package agents

import (
	"context"
	"fmt"
	"orion/internal/execution"
	"orion/internal/types"
	"time"
	"sync"
	"github.com/google/uuid"
)

type AutonomousAgent interface {
	Name() string
	Capabilities() []string
	ExecuteTask(ctx context.Context, task interface{}) (interface{}, error)
}

type AgentJob struct {
	IDVal    string
	Agent    AutonomousAgent
	Task     interface{}
	EventBus *types.EventBus
}
func (j *AgentJob) ID() string { return j.IDVal }
func (j *AgentJob) Type() string { return "AgentTask" }
func (j *AgentJob) Execute(ctx context.Context) error {
	j.EventBus.Publish(types.Event{Type: "agent.task_received", Payload: j.Agent.Name(), CreatedAt: time.Now()})
	_, err := j.Agent.ExecuteTask(ctx, j.Task)
	return err
}

type Registry struct { mu sync.RWMutex; agents map[string]AutonomousAgent }
func NewRegistry() *Registry { return &Registry{agents: make(map[string]AutonomousAgent)} }
func (r *Registry) RegisterAgent(a AutonomousAgent) { r.mu.Lock(); defer r.mu.Unlock(); r.agents[a.Name()] = a }
func (r *Registry) ListAgents() []AutonomousAgent { r.mu.RLock(); defer r.mu.RUnlock(); var l []AutonomousAgent; for _, a := range r.agents { l = append(l, a) }; return l }
func (r *Registry) GetAgentsByCapability(c string) []AutonomousAgent { r.mu.RLock(); defer r.mu.RUnlock(); var m []AutonomousAgent; for _, a := range r.agents { for _, cap := range a.Capabilities() { if cap == c { m = append(m, a); break } } }; return m }

type Supervisor struct { reg *Registry; eb *types.EventBus }
func NewSupervisor(r *Registry, eb *types.EventBus) *Supervisor { return &Supervisor{reg: r, eb: eb} }
func (s *Supervisor) StartAgents(ctx context.Context) error { return nil }
func (s *Supervisor) StopAgents(ctx context.Context) error { return nil }

type Dispatcher struct { reg *Registry; sch *execution.Scheduler; eb *types.EventBus }
func NewDispatcher(r *Registry, sch *execution.Scheduler, eb *types.EventBus) *Dispatcher { return &Dispatcher{reg: r, sch: sch, eb: eb} }
func (d *Dispatcher) Dispatch(ctx context.Context, cap string, task interface{}) error {
	m := d.reg.GetAgentsByCapability(cap)
	if len(m) == 0 { return fmt.Errorf("no agent for %s", cap) }
	d.sch.Schedule(&AgentJob{IDVal: uuid.New().String(), Agent: m[0], Task: task, EventBus: d.eb})
	return nil
}
