package runtime

import (
	"context"
	"fmt"
	"time"

	"orion/internal/cognition"
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
	id        string
	name      string
	agentType string
	eventBus  *EventBus
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

// ConversationAgent: Captures user intent and interacts via goals
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
	fmt.Printf("Agent [%s] (%s) starting...\n", a.name, a.agentType)

	a.eventBus.Subscribe("goal_created", func(event Event) {
		fmt.Printf("Agent [%s]: New goal received\n", a.name)
	})

	a.eventBus.Subscribe("goal_completed", func(event Event) {
		fmt.Printf("Agent [%s]: Goal completion recorded\n", a.name)
	})

	return nil
}

func (a *ConversationAgent) Stop(ctx context.Context) error {
	fmt.Printf("Agent [%s]: Stopping\n", a.name)
	return nil
}

// PlannerAgent: Breaks down goals into execution plans
type PlannerAgent struct {
	BaseAgent
	engine *cognition.Engine
}

func NewPlannerAgent(id, name string, eb *EventBus, engine *cognition.Engine) *PlannerAgent {
	return &PlannerAgent{
		BaseAgent: BaseAgent{
			id:        id,
			name:      name,
			agentType: "planner",
			eventBus:  eb,
		},
		engine: engine,
	}
}

func (a *PlannerAgent) Start(ctx context.Context) error {
	fmt.Printf("Agent [%s] (%s) starting...\n", a.name, a.agentType)

	a.eventBus.Subscribe("goal_created", func(event Event) {
		fmt.Printf("Agent [%s]: Planning execution steps for goal\n", a.name)

		// Trigger OODA-L cognition loop for the goal
		go func() {
			if err := a.engine.Process(ctx, "Analyze the goal: " + string(event.Payload)); err != nil {
				fmt.Printf("Agent [%s]: Cognition processing failed: %v\n", a.name, err)
			}
		}()
	})

	return nil
}

func (a *PlannerAgent) Stop(ctx context.Context) error {
	fmt.Printf("Agent [%s]: Stopping\n", a.name)
	return nil
}

// AnalysisAgent: Analyzes tool results and extracts insights
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
	fmt.Printf("Agent [%s] (%s) starting...\n", a.name, a.agentType)

	a.eventBus.Subscribe("tool_executed", func(event Event) {
		fmt.Printf("Agent [%s]: Analyzing tool execution result\n", a.name)
	})

	return nil
}

func (a *AnalysisAgent) Stop(ctx context.Context) error {
	fmt.Printf("Agent [%s]: Stopping\n", a.name)
	return nil
}

// CodeIndexerAgent: Parses and indexes source code symbols
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
	fmt.Printf("Agent [%s] (%s) starting...\n", a.name, a.agentType)
	return nil
}

func (a *CodeIndexerAgent) Stop(ctx context.Context) error {
	fmt.Printf("Agent [%s]: Stopping\n", a.name)
	return nil
}

// MemoryGardenerAgent: Deduplicates and consolidates knowledge
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
	fmt.Printf("Agent [%s] (%s) starting...\n", a.name, a.agentType)

	// Start periodic gardening
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Printf("Agent [%s]: Running knowledge consolidation\n", a.name)
			}
		}
	}()

	return nil
}

func (a *MemoryGardenerAgent) Stop(ctx context.Context) error {
	fmt.Printf("Agent [%s]: Stopping\n", a.name)
	return nil
}

// PatternDetectorAgent: Identifies recurring events and insights
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
	fmt.Printf("Agent [%s] (%s) starting...\n", a.name, a.agentType)

	a.eventBus.Subscribe("goal_completed", func(event Event) {
		fmt.Printf("Agent [%s]: Checking for patterns in completed goals\n", a.name)
	})

	return nil
}

func (a *PatternDetectorAgent) Stop(ctx context.Context) error {
	fmt.Printf("Agent [%s]: Stopping\n", a.name)
	return nil
}
