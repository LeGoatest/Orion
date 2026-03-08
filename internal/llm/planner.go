package llm

import (
	"context"
	"fmt"
)

// Planner implements the logic to convert LLM output into executable steps
type Planner struct {
	client *Client
}

func (p *Planner) GenerateNormalizedSteps(ctx context.Context, resp *Response) ([]string, error) {
	fmt.Println("LLM Planner: Normalizing response into execution steps")
	// Parsing logic to convert raw text/json into step list
	return resp.Steps, nil
}
