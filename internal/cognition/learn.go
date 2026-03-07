package cognition

import (
	"context"
	"fmt"
)

// Learn implementation: update knowledge from execution results
func (p *DefaultPipeline) Learn(ctx context.Context, result interface{}) error {
	fmt.Println("OODA-L: Phase: Learn")

	actResult, ok := result.(*ActResult)
	if !ok {
		return fmt.Errorf("invalid input type: expected *ActResult")
	}

	fmt.Printf("Learning from execution results for goal: %s\n", actResult.GoalID)

	// Learn logic:
	// 1. Convert act results into memory nodes
	// 2. Link results to the original goal node
	// 3. Trigger pattern detection if needed
	// 4. Update embeddings for retrieval engine

	// For bootstrap, just logging
	fmt.Printf("Knowledge graph updated for goal completion: %s\n", actResult.GoalID)
	return nil
}
