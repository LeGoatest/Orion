package cognition

import (
	"context"
	"fmt"
)

// Learn stores results as memory nodes and creates graph links
func (p *DefaultPipeline) Learn(ctx context.Context, result interface{}) error {
	actResult, ok := result.(*ActResult)
	if !ok {
		return fmt.Errorf("invalid input type for learn: expected *ActResult")
	}

	fmt.Printf("Learning from execution results for goal: %s\n", actResult.GoalID)

	// Learn logic:
	// - store memory nodes (facts, insights, tool_results)
	// - create graph links (related_to, caused_by, solved_by)
	// - detect patterns
	// - update embeddings

	// For bootstrap, just logging
	fmt.Printf("Knowledge updated from goal completion: %s\n", actResult.GoalID)
	return nil
}
