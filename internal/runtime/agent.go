package runtime

import (
	"context"
	"orion/internal/types"
)

// BaseAgent provides a common structure for agents
type BaseAgent struct {
	id        string
	name      string
	agentType string
	eventBus  *types.EventBus
}

func (a *BaseAgent) ID() string {
	return a.id
}

func (a *BaseAgent) AgentType() string {
	return a.agentType
}

func (a *BaseAgent) Start(ctx context.Context) error {
	return nil
}

func (a *BaseAgent) Stop(ctx context.Context) error {
	return nil
}
