package types

import (
	"sync"
	"time"
)

// Event represents an asynchronous system message
type Event struct {
	ID        string
	Type      string
	Payload   interface{}
	CreatedAt time.Time
}

// EventBus provides asynchronous communication with stage-aware routing
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Event),
	}
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

	// Direct type routing
	if subs, ok := eb.subscribers[event.Type]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
			}
		}
	}

	// Wildcard routing for loggers/supervisors
	if subs, ok := eb.subscribers["*"]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
			}
		}
	}
}
