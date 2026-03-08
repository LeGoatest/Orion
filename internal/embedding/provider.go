package embedding

import (
	"context"
)

// Provider defines the interface for generating embeddings
type Provider interface {
	Generate(ctx context.Context, text string) ([]float32, error)
	Name() string
}

// MockProvider is a stub for bootstrap
type MockProvider struct{}

func (p *MockProvider) Name() string { return "mock-provider" }

func (p *MockProvider) Generate(ctx context.Context, text string) ([]float32, error) {
	// Return a fixed dimension zero vector for now
	return make([]float32, 1536), nil
}
