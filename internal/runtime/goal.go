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
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	CurrentStage  GoalStage `json:"current_stage"`
	Status        string    `json:"status"`
	AssignedAgent string    `json:"assigned_agent"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GoalStore struct {
	db *sql.DB
}

func NewGoalStore(db *sql.DB) *GoalStore {
	return &GoalStore{db: db}
}

func (s *GoalStore) Save(g *Goal) error {
	query := `INSERT OR REPLACE INTO goals (id, description, current_stage, status, assigned_agent, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, g.ID, g.Description, string(g.CurrentStage), g.Status, g.AssignedAgent, g.CreatedAt, g.UpdatedAt)
	return err
}

func (s *GoalStore) UpdateStage(id string, next GoalStage, agent string) error {
	query := `UPDATE goals SET current_stage = ?, assigned_agent = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, string(next), agent, time.Now(), id)
	return err
}
