package agents

import (
	"context"
	"fmt"
	"orion/internal/runtime"
	"time"
)

type AutonomousAgent interface {
	Name() string
	Capabilities() []string
	ExecuteTask(context.Context, interface{}) (interface{}, error)
}

type AgentJob struct {
	IDVal  string
	GoalID string
	Stage  string
	Agent  AutonomousAgent
	Task   interface{}
	EB     *runtime.EventBus
}

func (j *AgentJob) ID() string {
	return j.IDVal
}

func (j *AgentJob) Type() string {
	return "AgentTask"
}

func (j *AgentJob) Execute(ctx context.Context) error {
	j.EB.Publish(runtime.Event{
		Type:      "agent.task_received",
		GoalID:    j.GoalID,
		Stage:     j.Stage,
		Payload:   j.Agent.Name(),
		CreatedAt: time.Now(),
	})

	r, err := j.Agent.ExecuteTask(ctx, j.Task)
	if err != nil {
		j.EB.Publish(runtime.Event{
			Type:   "agent.error",
			GoalID: j.GoalID,
			Stage:  j.Stage,
			Payload: map[string]string{
				"agent": j.Agent.Name(),
				"error": err.Error(),
			},
			CreatedAt: time.Now(),
		})
		return err
	}

	j.EB.Publish(runtime.Event{
		Type:   fmt.Sprintf("cognition.%s.completed", j.Stage),
		GoalID: j.GoalID,
		Stage:  j.Stage,
		Payload: map[string]interface{}{
			"goal_id": j.GoalID,
			"result":  r,
		},
		CreatedAt: time.Now(),
	})
	return nil
}
