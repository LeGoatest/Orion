package agents

import (
	"context"
	"fmt"
	"orion/internal/execution/scheduler"
	"orion/internal/types"
	"time"

	"github.com/google/uuid"
)

// Dispatcher assigns work to agents via the scheduler
type Dispatcher struct {
	registry  *Registry
	scheduler *scheduler.Scheduler
	eventBus  *types.EventBus
}

func NewDispatcher(registry *Registry, sch *scheduler.Scheduler, eb *types.EventBus) *Dispatcher {
	return &Dispatcher{
		registry:  registry,
		scheduler: sch,
		eventBus:  eb,
	}
}

// Dispatch matches a task to a capable agent and schedules its execution
func (d *Dispatcher) Dispatch(ctx context.Context, capability string, task interface{}) error {
	matches := d.registry.GetAgentsByCapability(capability)
	if len(matches) == 0 {
		return fmt.Errorf("no agent found for capability: %s", capability)
	}

	// Select best agent (simplified: first match)
	agent := matches[0]

	jobID := uuid.New().String()
	job := NewAgentJob(jobID, agent, task, d.eventBus)

	d.scheduler.Schedule(job)

	d.eventBus.Publish(types.Event{
		Type:      "agent.dispatch",
		Payload:   map[string]string{"agent": agent.Name(), "job_id": jobID},
		CreatedAt: time.Now(),
	})

	return nil
}
