package runtime

import (
	"context"
	"fmt"
	"time"
)

// Agent defines the interface for an Orion agent
type Agent interface {
	ID() string
	Name() string
	Type() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// BaseAgent provides a common structure for agents
type BaseAgent struct {
	id       string
	name     string
	agentType string
	eventBus *EventBus
}

func (a *BaseAgent) ID() string {
	return a.id
}

func (a *BaseAgent) Name() string {
	return a.name
}

func (a *BaseAgent) Type() string {
	return a.agentType
}

// ConversationAgent stub implementation
type ConversationAgent struct {
	BaseAgent
}

func NewConversationAgent(id, name string, eb *EventBus) *ConversationAgent {
	return &ConversationAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "conversation",
			eventBus:  eb,
		},
	}
}

func (a *ConversationAgent) Start(ctx context.Context) error {
	fmt.Printf("Starting %s (%s)\n", a.name, a.agentType)

	// Subscribe to relevant events
	a.eventBus.Subscribe("goal_created", func(event Event) {
		fmt.Printf("%s: Received goal_created event\n", a.name)
	})

	// Agent logic loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("%s: Stopping agent loop\n", a.name)
				return
			case <-time.After(1 * time.Minute):
				fmt.Printf("%s: Pulse check\n", a.name)
			}
		}
	}()

	return nil
}

func (a *ConversationAgent) Stop(ctx context.Context) error {
	fmt.Printf("Stopping %s\n", a.name)
	return nil
}

type AnalysisAgent struct {
	BaseAgent
}

func NewAnalysisAgent(id, name string, eb *EventBus) *AnalysisAgent {
	return &AnalysisAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "analysis",
			eventBus:  eb,
		},
	}
}

func (a *AnalysisAgent) Start(ctx context.Context) error {
	fmt.Printf("Starting %s (%s)\n", a.name, a.agentType)
	return nil
}

func (a *AnalysisAgent) Stop(ctx context.Context) error {
	fmt.Printf("Stopping %s\n", a.name)
	return nil
}

type MemoryGardenerAgent struct {
	BaseAgent
}

func NewMemoryGardenerAgent(id, name string, eb *EventBus) *MemoryGardenerAgent {
	return &MemoryGardenerAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "memory_gardener",
			eventBus:  eb,
		},
	}
}

func (a *MemoryGardenerAgent) Start(ctx context.Context) error {
	fmt.Printf("Starting %s (%s)\n", a.name, a.agentType)
	return nil
}

func (a *MemoryGardenerAgent) Stop(ctx context.Context) error {
	fmt.Printf("Stopping %s\n", a.name)
	return nil
}

type PatternDetectorAgent struct {
	BaseAgent
}

func NewPatternDetectorAgent(id, name string, eb *EventBus) *PatternDetectorAgent {
	return &PatternDetectorAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "pattern_detector",
			eventBus:  eb,
		},
	}
}

func (a *PatternDetectorAgent) Start(ctx context.Context) error {
	fmt.Printf("Starting %s (%s)\n", a.name, a.agentType)
	return nil
}

func (a *PatternDetectorAgent) Stop(ctx context.Context) error {
	fmt.Printf("Stopping %s\n", a.name)
	return nil
}

// Additional agent stubs can be added here
type PlannerAgent struct {
	BaseAgent
}

func NewPlannerAgent(id, name string, eb *EventBus) *PlannerAgent {
	return &PlannerAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "planner",
			eventBus:  eb,
		},
	}
}

func (a *PlannerAgent) Start(ctx context.Context) error {
	fmt.Printf("Starting %s (%s)\n", a.name, a.agentType)
	return nil
}

func (a *PlannerAgent) Stop(ctx context.Context) error {
	fmt.Printf("Stopping %s\n", a.name)
	return nil
}

type CodeIndexerAgent struct {
	BaseAgent
}

func NewCodeIndexerAgent(id, name string, eb *EventBus) *CodeIndexerAgent {
	return &CodeIndexerAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "code_indexer",
			eventBus:  eb,
		},
	}
}

func (a *CodeIndexerAgent) Start(ctx context.Context) error {
	fmt.Printf("Starting %s (%s)\n", a.name, a.agentType)
	return nil
}

func (a *CodeIndexerAgent) Stop(ctx context.Context) error {
	fmt.Printf("Stopping %s\n", a.name)
	return nil
}
