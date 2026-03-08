package agents

import (
	"context"
	"database/sql"
	"fmt"
	"orion/internal/execution"
	"orion/internal/types"
	"time"

	"github.com/google/uuid"
)

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

func (d *Dispatcher) Dispatch(ctx context.Context, db *sql.DB, workspaceID string, stage string, goalID string, task interface{}) error {
	// Mapping stage to agent capability
	cap := d.mapStageToCapability(stage)

	// Identify agents with required capability
	m := d.Registry.GetAgentsByCapability(cap)
	if len(m) == 0 {
		return fmt.Errorf("no agent for stage %s (capability %s)", stage, cap)
	}

	// Select the first available agent
	agent := m[0]
	jobID := uuid.New().String()

	// Transactional Job Creation in Workspace DB
	if db != nil {
		if err := d.persistJob(db, jobID, goalID, stage, agent.Name()); err != nil {
			return fmt.Errorf("failed to persist job: %w", err)
		}
	}

	// Construct the Job for the Scheduler
	job := &AgentJob{
		IDVal:       jobID,
		GoalID:      goalID,
		WorkspaceID: workspaceID,
		Stage:       stage,
		Agent:       agent,
		Task:        task,
		EventBus:    d.EventBus,
		WorkspaceDB: db,
	}

	// All execution flows through the Scheduler
	d.Scheduler.Schedule(job)

	// Notify system of job assignment
	d.EventBus.Publish(types.Event{
		Type:        "agent.task_assigned",
		GoalID:      goalID,
		WorkspaceID: workspaceID,
		Payload:     map[string]interface{}{"agent": agent.Name(), "stage": stage},
		CreatedAt:   time.Now(),
	})

	return nil
}

func (d *Dispatcher) mapStageToCapability(stage string) string {
	switch stage {
	case "ACT":
		return "tool_execution"
	case "OBSERVE":
		return "intake"
	case "PLAN":
		return "planning"
	default:
		return stage
	}
}

func (d *Dispatcher) persistJob(db *sql.DB, jobID, goalID, stage, agentName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Update Goal Stage in the same transaction
	_, err = tx.Exec("UPDATE goals SET current_stage = ?, updated_at = ? WHERE id = ?", stage, now, goalID)
	if err != nil {
		return err
	}

	// Insert new Job Record
	_, err = tx.Exec(`INSERT INTO jobs (id, goal_id, stage, assigned_agent, status, created_at, updated_at)
					  VALUES (?, ?, ?, ?, ?, ?, ?)`,
		jobID, goalID, stage, agentName, "pending", now, now)
	if err != nil {
		return err
	}

	return tx.Commit()
}
