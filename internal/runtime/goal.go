package runtime

import (
	"database/sql"
	"fmt"
	"time"
)

type GoalStage string

const (
	StageObserve      GoalStage = "OBSERVE"
	StageSymbolLookup GoalStage = "SYMBOL_LOOKUP"
	StagePatternMatch GoalStage = "PATTERN_MATCH"
	StageRetrieval    GoalStage = "RETRIEVAL"
	StagePlan         GoalStage = "PLAN"
	StageAct          GoalStage = "ACT"
	StageLearn        GoalStage = "LEARN"
	StageGarden       GoalStage = "GARDEN"
	StageCompleted    GoalStage = "COMPLETED"
	StageFailed       GoalStage = "FAILED"
)

type Goal struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	CurrentStage GoalStage `json:"current_stage"`
	Status       string    `json:"status"` // pending, active, completed, failed
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (g *Goal) AdvanceStage(next GoalStage) {
	fmt.Printf("Goal %s: %s -> %s\n", g.ID, g.CurrentStage, next)
	g.CurrentStage = next
	g.UpdatedAt = time.Now()
}

func (g *Goal) Persist(db *sql.DB) error {
	query := `INSERT INTO goals (id, description, current_stage, status, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)
			  ON CONFLICT(id) DO UPDATE SET
			  current_stage = excluded.current_stage,
			  status = excluded.status,
			  updated_at = excluded.updated_at`

	_, err := db.Exec(query, g.ID, g.Description, g.CurrentStage, g.Status, g.CreatedAt, g.UpdatedAt)
	return err
}

func LoadGoal(db *sql.DB, id string) (*Goal, error) {
	g := &Goal{}
	err := db.QueryRow("SELECT id, description, current_stage, status, created_at, updated_at FROM goals WHERE id = ?", id).
		Scan(&g.ID, &g.Description, &g.CurrentStage, &g.Status, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return g, nil
}
