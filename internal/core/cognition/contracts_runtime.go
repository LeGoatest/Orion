package cognition

import (
	"context"
	"fmt"
	"reflect"
)

// StageContract defines the expected types for a cognitive stage.
type StageContract struct {
	StageName  string
	InputType  reflect.Type
	OutputType reflect.Type
}

// Registry of known stage contracts.
var contracts = map[string]StageContract{
	"OBSERVE": {
		StageName:  "OBSERVE",
		InputType:  reflect.TypeOf(""),
		OutputType: reflect.TypeOf((*interface{})(nil)).Elem(), // Placeholder
	},
}

// ValidateContract ensures the input matches the stage's declared contract.
func ValidateContract(ctx context.Context, stage string, input interface{}) error {
	contract, ok := contracts[stage]
	if !ok {
		// If no contract is defined, we log a warning but allow for now to keep it extensible.
		fmt.Printf("Warning: No contract defined for stage %s\n", stage)
		return nil
	}

	inputType := reflect.TypeOf(input)
	if inputType != contract.InputType {
		return fmt.Errorf("contract violation for stage %s: expected %v, got %v", stage, contract.InputType, inputType)
	}

	return nil
}
