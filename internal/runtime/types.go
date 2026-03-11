package runtime

import (
	"orion/internal/types"
)

type Event = types.Event
type EventBus = types.EventBus

type BaseAgent struct {
	EventBus *types.EventBus
}
