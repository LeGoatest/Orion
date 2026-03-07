package cognition

import (
	"context"
)

type OODALPipeline interface {
	Observe(ctx context.Context, input interface{}) (interface{}, error)
	Orient(ctx context.Context, observation interface{}) (interface{}, error)
	Decide(ctx context.Context, orientation interface{}) (interface{}, error)
	Act(ctx context.Context, decision interface{}) (interface{}, error)
	Learn(ctx context.Context, result interface{}) error
	Garden(ctx context.Context, goalID string) error
}

type DefaultPipeline struct{}
