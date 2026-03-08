package types

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID        string
	GoalID    string
	WorkspaceID string
	Type      string
	Payload   interface{}
	CreatedAt time.Time
}

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]chan Event)}
}

func (eb *EventBus) Subscribe(eventType string) chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan Event, 100)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
	return ch
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	eb.publishTo(event, event.Type)
	eb.publishTo(event, "*")
}

func (eb *EventBus) publishTo(event Event, eventType string) {
	if subs, ok := eb.subscribers[eventType]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
				// Event sent successfully
			default:
				// Buffer full: drop oldest event and enqueue new one
				// We need to drain one from the channel to make room
				select {
				case dropped := <-ch:
					fmt.Printf("Warning: EventBus buffer full for type %s. Dropped oldest event: %v\n", eventType, dropped.Type)
				default:
					// Channel might have been drained by someone else in the meantime
				}

				// Try again to send the new event
				select {
				case ch <- event:
				default:
					// If it still fails, just skip (should be rare)
				}
			}
		}
	}
}
