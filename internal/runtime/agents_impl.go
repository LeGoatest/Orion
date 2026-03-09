package runtime

import (
	"context"
	"fmt"
)

type BaseAgent struct {
	EventBus *EventBus
}

type ConversationAgent struct {
	BaseAgent
}

func (a *ConversationAgent) Name() string {
	return "ConversationAgent"
}

func (a *ConversationAgent) Capabilities() []string {
	return []string{"OBSERVE"}
}

func (a *ConversationAgent) ExecuteTask(ctx context.Context, t interface{}) (interface{}, error) {
	fmt.Println("Agent [ConversationAgent]: Observe")
	return t, nil
}

type SymbolLookupAgent struct {
	BaseAgent
}

func (a *SymbolLookupAgent) Name() string {
	return "SymbolLookupAgent"
}

func (a *SymbolLookupAgent) Capabilities() []string {
	return []string{"SYMBOL_LOOKUP"}
}

func (a *SymbolLookupAgent) ExecuteTask(ctx context.Context, t interface{}) (interface{}, error) {
	fmt.Println("Agent [SymbolLookupAgent]: SymbolLookup")
	return t, nil
}

type PlannerAgent struct {
	BaseAgent
}

func (a *PlannerAgent) Name() string {
	return "PlannerAgent"
}

func (a *PlannerAgent) Capabilities() []string {
	return []string{"PLAN"}
}

func (a *PlannerAgent) ExecuteTask(ctx context.Context, t interface{}) (interface{}, error) {
	fmt.Println("Agent [PlannerAgent]: Plan")
	return t, nil
}

type RetrievalAgent struct {
	BaseAgent
}

func (a *RetrievalAgent) Name() string {
	return "RetrievalAgent"
}

func (a *RetrievalAgent) Capabilities() []string {
	return []string{"RETRIEVAL", "PATTERN_MATCH"}
}

func (a *RetrievalAgent) ExecuteTask(ctx context.Context, t interface{}) (interface{}, error) {
	fmt.Println("Agent [RetrievalAgent]: Retrieval/PatternMatch")
	return t, nil
}

type ToolExecutionAgent struct {
	BaseAgent
}

func (a *ToolExecutionAgent) Name() string {
	return "ToolExecutionAgent"
}

func (a *ToolExecutionAgent) Capabilities() []string {
	return []string{"ACT"}
}

func (a *ToolExecutionAgent) ExecuteTask(ctx context.Context, t interface{}) (interface{}, error) {
	fmt.Println("Agent [ToolExecutionAgent]: Act")
	return t, nil
}
