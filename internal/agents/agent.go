package agents

import (
	"context"
	"orion/internal/execution/worker"
	"orion/internal/types"
)

// AutonomousAgent defines the interface for agents that can execute tasks independently
type AutonomousAgent interface {
	Name() string
	Capabilities() []string
	Priority() int
	ExecuteTask(ctx context.Context, task interface{}) (interface{}, error)
}

// AgentJob wraps an agent execution as a worker pool job
type AgentJob struct {
	id        string
	agent     AutonomousAgent
	task      interface{}
	eventBus  *types.EventBus
}

func NewAgentJob(id string, agent AutonomousAgent, task interface{}, eb *types.EventBus) *AgentJob {
	return &AgentJob{
		id:       id,
		agent:    agent,
		task:     task,
		eventBus: eb,
	}
}

func (j *AgentJob) ID() string   { return j.id }
func (j *AgentJob) Type() string { return "AgentTask" }

func (j *AgentJob) Execute(ctx context.Context) error {
	j.eventBus.Publish(types.Event{Type: "agent.task_received", Payload: j.agent.Name()})

	_, err := j.agent.ExecuteTask(ctx, j.task)
	if err != nil {
		j.eventBus.Publish(types.Event{Type: "agent.error", Payload: err.Error()})
		return err
	}

	j.eventBus.Publish(types.Event{Type: "agent.task_completed", Payload: j.agent.Name()})
	return nil
}
