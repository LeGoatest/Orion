package prompt

import (
	"fmt"
	"strings"
)

// Section defines a distinct part of the prompt envelope
type Section string

const (
	SystemSection    Section = "SYSTEM"
	GoalSection      Section = "GOAL"
	FactsSection     Section = "FACTS"
	PatternsSection  Section = "PATTERNS"
	InsightsSection  Section = "INSIGHTS"
	CodeSection      Section = "CODE CONTEXT"
	CallGraphSection Section = "CALL GRAPH"
	UserTaskSection  Section = "USER TASK"
)

// Envelope represents a structured LLM input
type Envelope struct {
	Sections map[Section]string
	Order    []Section
}

func NewEnvelope() *Envelope {
	return &Envelope{
		Sections: make(map[Section]string),
		Order: []Section{
			SystemSection,
			GoalSection,
			FactsSection,
			PatternsSection,
			InsightsSection,
			CodeSection,
			CallGraphSection,
			UserTaskSection,
		},
	}
}

// Render combines all sections into a final prompt string
func (e *Envelope) Render() string {
	var builder strings.Builder
	for _, s := range e.Order {
		if content, ok := e.Sections[s]; ok && content != "" {
			builder.WriteString(fmt.Sprintf("### %s ###\n%s\n\n", string(s), content))
		}
	}
	return builder.String()
}
