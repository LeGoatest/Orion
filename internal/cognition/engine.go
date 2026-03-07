package cognition

import (
	"context"
)

// OODALPipeline represents the cognitive loop phases
type OODALPipeline interface {
	Observe(ctx context.Context, input interface{}) (interface{}, error)
	Orient(ctx context.Context, observation interface{}) (interface{}, error)
	Decide(ctx context.Context, orientation interface{}) (interface{}, error)
	Act(ctx context.Context, decision interface{}) (interface{}, error)
	Learn(ctx context.Context, result interface{}) error
}

// Engine drives the OODA-L cognitive loop
type Engine struct {
	pipeline OODALPipeline
}

// NewEngine creates a new Engine
func NewEngine(pipeline OODALPipeline) *Engine {
	return &Engine{
		pipeline: pipeline,
	}
}

// Process runs a single iteration of the OODA-L loop
func (e *Engine) Process(ctx context.Context, input interface{}) error {
	// Observe
	observation, err := e.pipeline.Observe(ctx, input)
	if err != nil {
		return err
	}

	// Orient
	orientation, err := e.pipeline.Orient(ctx, observation)
	if err != nil {
		return err
	}

	// Decide
	decision, err := e.pipeline.Decide(ctx, orientation)
	if err != nil {
		return err
	}

	// Act
	result, err := e.pipeline.Act(ctx, decision)
	if err != nil {
		return err
	}

	// Learn
	err = e.pipeline.Learn(ctx, result)
	if err != nil {
		return err
	}

	return nil
}

// DefaultPipeline provides a basic implementation of the OODA-L loop
type DefaultPipeline struct{}
