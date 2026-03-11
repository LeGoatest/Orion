package cognition

import (
	"fmt"
)

func (ce *Engine) Decide(orientation interface{}) interface{} {
	fmt.Println("Cognition: Phase [Decide]")
	// Construct execution plan, select tools, generate plan steps
	return orientation
}
