package llm

import (
	"context"
	"fmt"
	"orion/internal/prompt"
	"orion/internal/types"
	"time"
)

// Client handles interaction with the LLM provider and emits events
type Client struct {
	provider Provider
	eb       *types.EventBus
}

func NewClient(provider Provider, eb *types.EventBus) *Client {
	return &Client{provider: provider, eb: eb}
}

func (c *Client) RequestPlan(ctx context.Context, envelope *prompt.Envelope) (*Response, error) {
	c.eb.Publish(types.Event{Type: "llm.requested", CreatedAt: time.Now()})

	resp, err := c.provider.Plan(ctx, envelope)
	if err != nil {
		c.eb.Publish(types.Event{
			Type:    "llm.failed",
			Payload: map[string]string{"error": err.Error()},
			CreatedAt: time.Now(),
		})
		return nil, fmt.Errorf("llm plan request failed: %w", err)
	}

	c.eb.Publish(types.Event{Type: "llm.completed", CreatedAt: time.Now()})
	return resp, nil
}
