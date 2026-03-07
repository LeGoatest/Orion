package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Goal status constants
const (
	GoalPending   = "pending"
	GoalActive    = "active"
	GoalCompleted = "completed"
	GoalFailed    = "failed"
)

// Goal represents a unit of work within a workspace
type Goal struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	ParentID    *string   `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GoalRuntime manages the goal lifecycle for a workspace
type GoalRuntime struct {
	db       *sql.DB
	eventBus *EventBus
}

// NewGoalRuntime creates a new GoalRuntime
func NewGoalRuntime(db *sql.DB, eventBus *EventBus) *GoalRuntime {
	return &GoalRuntime{
		db:       db,
		eventBus: eventBus,
	}
}

// CreateGoal initializes a new goal in the workspace database
func (gr *GoalRuntime) CreateGoal(ctx context.Context, description string, priority int, parentID *string) (*Goal, error) {
	id := uuid.New().String()
	goal := &Goal{
		ID:          id,
		Description: description,
		Status:      GoalPending,
		Priority:    priority,
		ParentID:    parentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO goals (id, description, status, priority, parent_id, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := gr.db.ExecContext(ctx, query, goal.ID, goal.Description, goal.Status, goal.Priority, goal.ParentID, goal.CreatedAt, goal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	// Publish goal_created event
	gr.eventBus.Publish(ctx, Event{
		Type:      "goal_created",
		Source:    "goal_runtime",
		Payload:   nil, // Could be marshaled goal
		Timestamp: time.Now(),
	})

	return goal, nil
}

// UpdateGoalStatus updates the status of an existing goal
func (gr *GoalRuntime) UpdateGoalStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE goals SET status = ?, updated_at = ? WHERE id = ?`
	_, err := gr.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update goal status: %w", err)
	}

	if status == GoalCompleted {
		gr.eventBus.Publish(ctx, Event{
			Type:      "goal_completed",
			Source:    "goal_runtime",
			Payload:   nil,
			Timestamp: time.Now(),
		})
	}

	return nil
}

// ListGoals retrieves all goals for the workspace
func (gr *GoalRuntime) ListGoals(ctx context.Context) ([]*Goal, error) {
	query := `SELECT id, description, status, priority, parent_id, created_at, updated_at FROM goals`
	rows, err := gr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list goals: %w", err)
	}
	defer rows.Close()

	var goals []*Goal
	for rows.Next() {
		var g Goal
		err := rows.Scan(&g.ID, &g.Description, &g.Status, &g.Priority, &g.ParentID, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, &g)
	}
	return goals, nil
}
