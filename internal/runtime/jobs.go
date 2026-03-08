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
	CompletedAt   *time.Time `json:"completed_at"`
}

func (j *JobRecord) Persist(db *sql.DB) error {
	query := `INSERT INTO jobs (id, goal_id, stage, assigned_agent, status, retry_count, created_at, completed_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			  ON CONFLICT(id) DO UPDATE SET
			  status = excluded.status,
			  retry_count = excluded.retry_count,
			  completed_at = excluded.completed_at`

	_, err := db.Exec(query, j.ID, j.GoalID, j.Stage, j.AssignedAgent, j.Status, j.RetryCount, j.CreatedAt, j.CompletedAt)
	return err
}

func LoadJob(db *sql.DB, id string) (*JobRecord, error) {
	j := &JobRecord{}
	err := db.QueryRow("SELECT id, goal_id, stage, assigned_agent, status, retry_count, created_at, completed_at FROM jobs WHERE id = ?", id).
		Scan(&j.ID, &j.GoalID, &j.Stage, &j.AssignedAgent, &j.Status, &j.RetryCount, &j.CreatedAt, &j.CompletedAt)
	if err != nil {
		return nil, err
	}
	return j, nil
}
