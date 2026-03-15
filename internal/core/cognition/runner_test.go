package cognition

import (
	"context"
	"testing"
)

func TestStageRunner(t *testing.T) {
	reg := NewExecutorRegistry()
	runner := NewStageRunner(reg)

	// Test OBSERVE stage with correct input
	_, err := runner.RunStage(context.Background(), "OBSERVE", "ALGO_OBSERVE_DETERMINISTIC", "User intent")
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// Test unknown algorithm
	_, err = runner.RunStage(context.Background(), "OBSERVE", "UNKNOWN_ALGO", "input")
	if err == nil {
		t.Error("Expected error for unknown algorithm, got nil")
	}
}
