package runtime

import (
	"database/sql"
	"time"
)

type JobRecord struct {
	ID            string    `json:"id"`
	GoalID        string    `json:"goal_id"`
	Stage         string    `json:"stage"`
	AssignedAgent string    `json:"assigned_agent"`
	Status        string    `json:"status"` // pending, running, completed, failed
	RetryCount    int       `json:"retry_count"`
	CreatedAt     time.Time `json:"created_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

type JobStore struct {
	db *sql.DB
}

func NewJobStore(db *sql.DB) *JobStore {
	return &JobStore{db: db}
}

func (s *JobStore) Create(j *JobRecord) error {
	query := `INSERT INTO jobs (id, goal_id, stage, assigned_agent, status, retry_count, created_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, j.ID, j.GoalID, j.Stage, j.AssignedAgent, j.Status, j.RetryCount, j.CreatedAt)
	return err
}

func (s *JobStore) Complete(id string) error {
	query := `UPDATE jobs SET status = 'completed', finished_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, time.Now(), id)
	return err
}
