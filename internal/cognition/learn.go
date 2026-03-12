package cognition

import (
	"fmt"
	"time"
)

func (ce *Engine) Learn(outcome *OutcomeModel) {
	fmt.Printf("Cognition: Phase [Learn] - Persisting outcome for goal: %s\n", outcome.GoalID)

	// Record outcome tied to strategy and situational context
	// This will be used in future Orient phases to improve strategy choice
	outcome.Timestamp = time.Now()
	outcome.OutcomeSummary = fmt.Sprintf("Strategy %s result: %v",
		outcome.DecisionModel.StrategyChoice, outcome.ActionResult)

	// Persistence logic would go here
	fmt.Printf("Outcome: %s (Success: %v)\n", outcome.OutcomeSummary, outcome.Success)
}
