package llm

import (
	"context"
	"orion/internal/prompt"
)

// Provider defines the interface for an LLM planning source
type Provider interface {
	Plan(ctx context.Context, envelope *prompt.Envelope) (*Response, error)
	Name() string
}

// StubProvider is used when no real LLM is connected
type StubProvider struct{}

func (p *StubProvider) Name() string { return "stub-llm" }

func (p *StubProvider) Plan(ctx context.Context, envelope *prompt.Envelope) (*Response, error) {
	return &Response{
		Steps: []string{"step 1", "step 2"},
		Raw:   "stub response",
	}, nil
}
