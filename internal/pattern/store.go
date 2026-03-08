package pattern

import (
	"context"
	"database/sql"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// MatchTrigger searches for a pattern trigger that matches the input goal
func (s *Store) MatchTrigger(ctx context.Context, goal string) (*Pattern, bool) {
	query := `SELECT id, trigger, solution_steps, confidence, usage_count, state FROM patterns WHERE ? LIKE '%' || trigger || '%'`
	row := s.db.QueryRowContext(ctx, query, goal)

	var p Pattern
	var stepsJSON string
	err := row.Scan(&p.ID, &p.Trigger, &stepsJSON, &p.Confidence, &p.UsageCount, &p.State)
	if err != nil {
		return nil, false
	}

	// Unmarshal solution steps
	steps := strings.Split(stepsJSON, ",")
	p.SolutionSteps = steps

	return &p, true
}

// SavePattern persists a new pattern to the database
func (s *Store) SavePattern(ctx context.Context, p *Pattern) error {
	stepsJSON := strings.Join(p.SolutionSteps, ",")
	query := `INSERT INTO patterns (id, trigger, solution_steps, confidence, usage_count, state)
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, p.ID, p.Trigger, stepsJSON, p.Confidence, p.UsageCount, p.State)
	return err
}

// IncrementUsage updates statistics for a matched pattern
func (s *Store) IncrementUsage(ctx context.Context, id string) error {
	query := `UPDATE patterns SET usage_count = usage_count + 1, confidence = MIN(1.0, confidence + 0.05) WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
