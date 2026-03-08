# GAP REPORT - Agent-Driven Cognition Pipeline

## Existing Systems
- Multi-agent foundation (Registry, Supervisor, Dispatcher)
- Centralized cognition engine (cognition.go)
- Event bus
- SQLite storage

## Missing/Partial Logic
- **Agent Stage Ownership**: Agents are not yet explicitly assigned to cognition stages.
- **Dispatcher Stage Routing**: Dispatcher doesn't map stages to capabilities or create durable jobs.
- **Goal Stage Persistence**: Goals do not track their current OODA-L stage in the database.
- **Event-Based Handoff**: The engine triggers stages linearly rather than responding to bus events for handoff.

## Architecture Conflicts
- Transition from centralized function calls to distributed agent jobs requires careful state management to ensure determinism.
