# GAP REPORT - Agent Coordination in the Cognition Pipeline

## Existing Systems
- Multi-agent runtime foundation (internal/agents/)
- OODA-L cognition engine (internal/cognition/)
- Event bus (internal/types/event_bus.go)
- Goal and job models (internal/runtime/goal/)

## Missing/Partial Features
- **Agent Stage Ownership**: Agents are not yet assigned specific responsibilities within the cognition pipeline.
- **Event-Driven Handoff**: No events to signal completion of one stage and trigger the next via dispatcher.
- **Stage Tracking**: Goals don't explicitly track which cognition stage they are currently in.
- **ToolExecutorAgent**: Missing an agent specialized in the ACT stage execution.
- **Agent Lifecycle Events**: Mandatory events (agent.started, agent.failed, etc.) are partially implemented.

## Architecture Conflicts
- Cognition loop is currently a single linear function; it must be refactored into a sequence of dispatchable, event-driven agent tasks.
