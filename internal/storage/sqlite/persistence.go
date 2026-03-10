package sqlite

import (
	"context"
	"orion/ent"
	"orion/ent/goal"
	"orion/ent/job"
	"time"
)

type Persistence struct {
	client *ent.Client
}

func NewPersistence(client *ent.Client) *Persistence {
	return &Persistence{client: client}
}

func (p *Persistence) TransitionStage(ctx context.Context, goalID string, nextStage string, agentName string) error {
	tx, err := p.client.Tx(ctx)
	if err != nil {
		return err
	}

	// Update goal
	err = tx.Goal.Update().
		Where(goal.ID(goalID)).
		SetCurrentStage(nextStage).
		SetAssignedAgent(agentName).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create/Update job record
	_, err = tx.Job.Create().
		SetID(goalID + "-" + nextStage).
		SetGoalID(goalID).
		SetStage(nextStage).
		SetAssignedAgent(agentName).
		SetStatus("RUNNING").
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
