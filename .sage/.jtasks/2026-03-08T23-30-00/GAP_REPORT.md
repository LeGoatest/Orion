# GAP REPORT - Agent-Driven Cognition Pipeline

## Existing Systems
- Multi-agent foundation (Registry, Supervisor, Dispatcher)
- OODA-L cognition logic (centralized)
- Event bus
- Goal persistence

## Missing Components
- **Agent Stage Ownership**: Logic to explicitly assign and hand off cognition stages to agents.
- **Event-Based Handoff**: The engine currently triggers stages linearly; it needs to respond to completion events to trigger the next dispatch.
- **Stage Tracking**: Goals don't yet record their progress through the specific cognition stages (SymbolLookup, PatternMatch, etc.).
- **Autonomous Background Loops**: Background agents (Indexer, Gardener) need autonomous control within the multi-agent framework.

## Architecture Conflicts
- Transition from a "function-call" pipeline to an "event-dispatch" pipeline requires careful handling of shared context (ContextBundle).
