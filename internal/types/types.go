package types

import (
	"context"
	"fmt"
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

	publishTo := func(subs []chan Event) {
		for _, ch := range subs {
			select {
			case ch <- e:
			default:
				select {
				case <-ch:
				default:
				}
				select {
				case ch <- e:
					fmt.Printf("Warning: Event bus slow consumer for type %s. Dropped oldest event.\n", e.Type)
				default:
				}
			}
		}
	}

	if s, ok := eb.subs[e.Type]; ok {
		publishTo(s)
	}
	if s, ok := eb.subs["*"]; ok {
		publishTo(s)
	}
}

type Job func(ctx context.Context)
