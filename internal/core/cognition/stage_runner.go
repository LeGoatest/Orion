package cognition

import (
	"context"
	"fmt"
)

// StageRunner orchestrates the execution of a cognitive stage using the registry.
type StageRunner struct {
	registry *ExecutorRegistry
}

func NewStageRunner(r *ExecutorRegistry) *StageRunner {
	return &StageRunner{registry: r}
}

// RunStage resolves the algorithm, validates the contract, and executes the logic.
func (sr *StageRunner) RunStage(ctx context.Context, stageName string, algoID string, input interface{}) (interface{}, error) {
	// 1. Contract Validation
	if err := ValidateContract(ctx, stageName, input); err != nil {
		return nil, fmt.Errorf("stage %s: %w", stageName, err)
	}

	// 2. Algorithm Resolution
	algo, err := sr.registry.Get(algoID)
	if err != nil {
		return nil, fmt.Errorf("stage %s: %w", stageName, err)
	}

	// 3. Deterministic Execution
	fmt.Printf("Executing stage %s using algorithm %s\n", stageName, algoID)
	output, err := algo.Execute(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("stage %s execution failed: %w", stageName, err)
	}

	return output, nil
}
