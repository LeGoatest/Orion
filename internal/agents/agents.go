package agents

import (
	"context"
	"fmt"
	"orion/internal/execution"
	"orion/internal/types"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AutonomousAgent interface {
	Name() string
	Capabilities() []string
	ExecuteTask(ctx context.Context, task interface{}) (interface{}, error)
}

type AgentJob struct {
	IDVal      string
	GoalID     string
	Stage      string
	Agent      AutonomousAgent
	Task       interface{}
	EventBus   *types.EventBus
	RetryCount int
}

func (j *AgentJob) ID() string   { return j.IDVal }
func (j *AgentJob) Type() string { return "AgentTask" }

func (j *AgentJob) Execute(ctx context.Context) error {
	j.EventBus.Publish(types.Event{
		Type:      "agent.task_received",
		Payload:   map[string]string{"agent": j.Agent.Name(), "stage": j.Stage, "goal_id": j.GoalID},
		CreatedAt: time.Now(),
	})

	result, err := j.Agent.ExecuteTask(ctx, j.Task)
	if err != nil {
		j.EventBus.Publish(types.Event{
			Type:      "agent.error",
			Payload:   map[string]string{"agent": j.Agent.Name(), "error": err.Error(), "goal_id": j.GoalID},
			CreatedAt: time.Now(),
		})
		return err
	}

	// Publish stage completion event
	j.EventBus.Publish(types.Event{
		Type:      fmt.Sprintf("cognition.%s.completed", j.Stage),
		Payload:   map[string]interface{}{"goal_id": j.GoalID, "result": result},
		CreatedAt: time.Now(),
	})

	j.EventBus.Publish(types.Event{
		Type:      "agent.task_completed",
		Payload:   map[string]string{"agent": j.Agent.Name(), "stage": j.Stage, "goal_id": j.GoalID},
		CreatedAt: time.Now(),
	})

	return nil
}

type Registry struct {
	mu     sync.RWMutex
	agents map[string]AutonomousAgent
}

func NewRegistry() *Registry {
	return &Registry{agents: make(map[string]AutonomousAgent)}
}

func (r *Registry) RegisterAgent(a AutonomousAgent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[a.Name()] = a
}

func (r *Registry) ListAgents() []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var l []AutonomousAgent
	for _, a := range r.agents {
		l = append(l, a)
	}
	return l
}

func (r *Registry) GetAgentsByCapability(c string) []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var m []AutonomousAgent
	for _, a := range r.agents {
		for _, cap := range a.Capabilities() {
			if cap == c {
				m = append(m, a)
				break
			}
		}
	}
	return m
}

type Supervisor struct {
	Registry *Registry
	EventBus *types.EventBus
}

func NewSupervisor(r *Registry, eb *types.EventBus) *Supervisor {
	return &Supervisor{Registry: r, EventBus: eb}
}

func (s *Supervisor) StartAgents(ctx context.Context) error {
	for _, a := range s.Registry.ListAgents() {
		s.EventBus.Publish(types.Event{Type: "agent.started", Payload: a.Name(), CreatedAt: time.Now()})
	}
	return nil
}

func (s *Supervisor) StopAgents(ctx context.Context) error {
	for _, a := range s.Registry.ListAgents() {
		s.EventBus.Publish(types.Event{Type: "agent.stopped", Payload: a.Name(), CreatedAt: time.Now()})
	}
	return nil
}

type Dispatcher struct {
	Registry  *Registry
	Scheduler *execution.Scheduler
	EventBus  *types.EventBus
}

func NewDispatcher(r *Registry, s *execution.Scheduler, eb *types.EventBus) *Dispatcher {
	return &Dispatcher{
		Registry:  r,
		Scheduler: s,
		EventBus:  eb,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, stage string, goalID string, task interface{}) error {
	cap := stage
	if stage == "ACT" {
		cap = "tool_execution"
	}

	m := d.Registry.GetAgentsByCapability(cap)
	if len(m) == 0 {
		return fmt.Errorf("no agent for stage %s (capability %s)", stage, cap)
	}

	agent := m[0]
	job := &AgentJob{
		IDVal:    uuid.New().String(),
		GoalID:   goalID,
		Stage:    stage,
		Agent:    agent,
		Task:     task,
		EventBus: d.EventBus,
	}

	d.Scheduler.Schedule(job)

	d.EventBus.Publish(types.Event{
		Type:      "agent.task_assigned",
		Payload:   map[string]string{"agent": agent.Name(), "stage": stage, "goal_id": goalID},
		CreatedAt: time.Now(),
	})

	return nil
}
