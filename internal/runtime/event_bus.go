package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Event represents a system event
type Event struct {
	ID        int64           `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// EventSubscriber is a function that processes events
type EventSubscriber func(Event)

// EventBus manages in-memory event distribution and persistence
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventSubscriber
	db          *sql.DB // Global DB for persistence
}

// NewEventBus creates a new EventBus
func NewEventBus(db *sql.DB) *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventSubscriber),
		db:          db,
	}
}

// Subscribe adds a subscriber for a specific event type or "*" for all
func (eb *EventBus) Subscribe(eventType string, subscriber EventSubscriber) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
}

// Publish distributes an event to subscribers and persists if required
func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	eb.mu.RLock()
	subs := eb.subscribers[event.Type]
	allSubs := eb.subscribers["*"]
	eb.mu.RUnlock()

	// Distributed asynchronously to avoid blocking the publisher
	// Each subscriber runs in its own goroutine for isolation and robustness
	for _, sub := range subs {
		go safeInvoke(sub, event)
	}
	for _, sub := range allSubs {
		go safeInvoke(sub, event)
	}

	// Persist critical events if needed
	if shouldPersist(event.Type) {
		if err := eb.persist(ctx, event); err != nil {
			return fmt.Errorf("failed to persist event: %w", err)
		}
	}

	return nil
}

func (eb *EventBus) persist(ctx context.Context, event Event) error {
	if eb.db == nil {
		return nil
	}

	query := `INSERT INTO event_log (event_type, source, payload, timestamp) VALUES (?, ?, ?, ?)`
	_, err := eb.db.ExecContext(ctx, query, event.Type, event.Source, string(event.Payload), event.Timestamp)
	return err
}

func safeInvoke(sub EventSubscriber, event Event) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in event subscriber: %v\n", r)
		}
	}()
	sub(event)
}

func shouldPersist(eventType string) bool {
	// Persist critical events for audit and learning
	switch eventType {
	case "goal_created", "goal_completed", "tool_executed", "pattern_detected":
		return true
	default:
		return false
	}
}
