package agents

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"orion/internal/execution"
	"orion/internal/types"
	"sync"
	"time"
)

type AutonomousAgent interface {
	Name() string
	Capabilities() []string
	ExecuteTask(context.Context, interface{}) (interface{}, error)
}

type AgentJob struct {
	IDVal  string
	GoalID string
	Stage  string
	Agent  AutonomousAgent
	Task   interface{}
	EB     *types.EventBus
}

func (j *AgentJob) ID() string   { return j.IDVal }
func (j *AgentJob) Type() string { return "AgentTask" }

func (j *AgentJob) Execute(ctx context.Context) error {
	j.EB.Publish(types.Event{
		Type:      "agent.task_received",
		GoalID:    j.GoalID,
		Stage:     j.Stage,
		Payload:   j.Agent.Name(),
		CreatedAt: time.Now(),
	})

	r, err := j.Agent.ExecuteTask(ctx, j.Task)
	if err != nil {
		j.EB.Publish(types.Event{
			Type:   "agent.error",
			GoalID: j.GoalID,
			Stage:  j.Stage,
			Payload: map[string]string{
				"agent": j.Agent.Name(),
				"error": err.Error(),
			},
			CreatedAt: time.Now(),
		})
		return err
	}

	j.EB.Publish(types.Event{
		Type:   fmt.Sprintf("cognition.%s.completed", j.Stage),
		GoalID: j.GoalID,
		Stage:  j.Stage,
		Payload: map[string]interface{}{
			"goal_id": j.GoalID,
			"result":  r,
		},
		CreatedAt: time.Now(),
	})
	return nil
}

type Registry struct {
	mu sync.RWMutex
	m  map[string]AutonomousAgent
}

func NewRegistry() *Registry {
	return &Registry{m: make(map[string]AutonomousAgent)}
}

func (r *Registry) RegisterAgent(a AutonomousAgent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[a.Name()] = a
}

func (r *Registry) GetByCap(c string) []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []AutonomousAgent
	for _, a := range r.m {
		for _, cap := range a.Capabilities() {
			if cap == c {
				res = append(res, a)
				break
			}
		}
	}
	return res
}

func (r *Registry) List() []AutonomousAgent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []AutonomousAgent
	for _, a := range r.m {
		res = append(res, a)
	}
	return res
}

type Supervisor struct {
	Reg    *Registry
	EB     *types.EventBus
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	status map[string]time.Time
}

func NewSupervisor(r *Registry, eb *types.EventBus) *Supervisor {
	return &Supervisor{
		Reg:    r,
		EB:     eb,
		status: make(map[string]time.Time),
	}
}

func (s *Supervisor) StartAgents(ctx context.Context) error {
	s.mu.Lock()
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.mu.Unlock()

	for _, a := range s.Reg.List() {
		s.startAgent(a)
	}

	go s.monitorLoop()
	return nil
}

func (s *Supervisor) startAgent(a AutonomousAgent) {
	s.EB.Publish(types.Event{
		Type:      "agent.started",
		Payload:   a.Name(),
		CreatedAt: time.Now(),
	})
	s.mu.Lock()
	s.status[a.Name()] = time.Now()
	s.mu.Unlock()
}

func (s *Supervisor) monitorLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			for _, a := range s.Reg.List() {
				lastSeen, exists := s.status[a.Name()]
				if !exists || now.Sub(lastSeen) > 15*time.Second {
					s.startAgent(a)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Supervisor) StopAgents() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Supervisor) Heartbeat(agentName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[agentName] = time.Now()
}

type Dispatcher struct {
	Reg *Registry
	Sch *execution.Scheduler
	EB  *types.EventBus
}

func NewDispatcher(r *Registry, s *execution.Scheduler, eb *types.EventBus) *Dispatcher {
	return &Dispatcher{Reg: r, Sch: s, EB: eb}
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
