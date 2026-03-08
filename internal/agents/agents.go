package agents

import (
	"context"
	"database/sql"
	"fmt"
	"orion/internal/types"
	"sync"
	"time"
)

type AutonomousAgent interface {
	Name() string
	Capabilities() []string
	ExecuteTask(ctx context.Context, task interface{}) (interface{}, error)
}

type AgentJob struct {
	IDVal       string
	GoalID      string
	WorkspaceID string
	Stage       string
	Agent       AutonomousAgent
	Task        interface{}
	EventBus    *types.EventBus
	RetryCount  int
	WorkspaceDB *sql.DB
}

func (j *AgentJob) ID() string   { return j.IDVal }
func (j *AgentJob) Type() string { return "AgentTask" }

func (j *AgentJob) Execute(ctx context.Context) error {
	j.EventBus.Publish(types.Event{
		Type:        "agent.task_received",
		GoalID:      j.GoalID,
		WorkspaceID: j.WorkspaceID,
		Payload:     map[string]interface{}{"agent": j.Agent.Name(), "stage": j.Stage},
		CreatedAt:   time.Now(),
	})

	// Mark job as running in DB
	if j.WorkspaceDB != nil {
		_, _ = j.WorkspaceDB.Exec("UPDATE jobs SET status = ?, updated_at = ? WHERE id = ?", "running", time.Now(), j.IDVal)
	}

	result, err := j.Agent.ExecuteTask(ctx, j.Task)
	if err != nil {
		j.EventBus.Publish(types.Event{
			Type:        "agent.error",
			GoalID:      j.GoalID,
			WorkspaceID: j.WorkspaceID,
			Payload:     map[string]interface{}{"agent": j.Agent.Name(), "error": err.Error()},
			CreatedAt:   time.Now(),
		})

		if j.WorkspaceDB != nil {
			_, _ = j.WorkspaceDB.Exec("UPDATE jobs SET status = ?, updated_at = ? WHERE id = ?", "failed", time.Now(), j.IDVal)
		}
		return err
	}

	// Publish stage completion event
	j.EventBus.Publish(types.Event{
		Type:        fmt.Sprintf("cognition.%s.completed", j.Stage),
		GoalID:      j.GoalID,
		WorkspaceID: j.WorkspaceID,
		Payload:     map[string]interface{}{"result": result},
		CreatedAt:   time.Now(),
	})

	j.EventBus.Publish(types.Event{
		Type:        "agent.task_completed",
		GoalID:      j.GoalID,
		WorkspaceID: j.WorkspaceID,
		Payload:     map[string]interface{}{"agent": j.Agent.Name(), "stage": j.Stage},
		CreatedAt:   time.Now(),
	})

	// Mark job as completed in DB
	if j.WorkspaceDB != nil {
		now := time.Now()
		_, _ = j.WorkspaceDB.Exec("UPDATE jobs SET status = ?, completed_at = ?, updated_at = ? WHERE id = ?", "completed", now, now, j.IDVal)
	}

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
