package runtime

import "time"

type JobRecord struct {
	ID            string    `json:"id"`
	GoalID        string    `json:"goal_id"`
	Stage         string    `json:"stage"`
	AssignedAgent string    `json:"assigned_agent"`
	Status        string    `json:"status"` // pending, running, completed, failed
	RetryCount    int       `json:"retry_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
