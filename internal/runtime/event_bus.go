package runtime

import (
	"sync"
	"time"
)

type Event struct {
	ID        string
	GoalID    string
	Stage     string
	Type      string
	Payload   interface{}
	CreatedAt time.Time
}

type EventBus struct {
	mu   sync.RWMutex
	subs map[string][]chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subs: make(map[string][]chan Event),
	}
}

func (eb *EventBus) Subscribe(t string) chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan Event, 100)
	eb.subs[t] = append(eb.subs[t], ch)
	return ch
}

func (eb *EventBus) Publish(e Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if s, ok := eb.subs[e.Type]; ok {
		for _, ch := range s {
			select {
			case ch <- e:
			default:
			}
		}
	}
	if s, ok := eb.subs["*"]; ok {
		for _, ch := range s {
			select {
			case ch <- e:
			default:
			}
		}
	}
}
