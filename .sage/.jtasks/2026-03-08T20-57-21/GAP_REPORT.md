# GAP REPORT — Orion Runtime Production Hardening

## 1. Goal + Job Persistence
- **Current State**: `Goal` and `JobRecord` are in-memory structs. Stage transitions in `internal/runtime/goal.go` only update fields and print to stdout.
- **Weakness**: System restarts lose all progress. There is no SQLite schema for goals or jobs.
- **Gap**: Missing transactional persistence in `workspace.db` for stage transitions.

## 2. Event Bus Reliability
- **Current State**: Simple in-memory channel-based pub/sub with a fixed buffer size of 100.
- **Weakness**: Slow consumers can cause the `Publish` method to block or drop events silently (current code uses `default` in `select` which drops events if buffer full, but no notification/logging).
- **Gap**: Need bounded buffers that drop oldest events and log warnings instead of silent drops.

## 3. Agent Supervisor Lifecycle
- **Current State**: `Supervisor` just publishes `agent.started` events at boot.
- **Weakness**: No active monitoring. If an agent goroutine crashes or hangs, the system is unaware. No restart policy.
- **Gap**: Missing heartbeat monitoring, crash detection, and automated restart logic.

## 4. Dispatcher Stability
- **Current State**: `Dispatcher` directly creates `AgentJob` and schedules it. It also retrieves the agent from the registry.
- **Weakness**: Responsibilities are mostly clean, but it lacks retry logic and tight integration with the persistent job records.
- **Gap**: Ensure dispatcher only maps stage to capability and creates/submits jobs via the scheduler, with full integration into the persistence layer.

## 5. Retrieval Pipeline Maturity
- **Current State**: `RetrievalEngine` has placeholders for vector search and graph expansion. `ContextBundle` is minimal.
- **Weakness**: No hybrid scoring implementation. Context assembly is mostly simulated.
- **Gap**: Missing weighted hybrid scoring algorithm and full assembly of `ContextBundle` (symbols + patterns + vector + graph).

## 6. Bootstrap Assumptions
- **Current State**: `main.go` has a hardcoded test goal "goal-coordinated-001".
- **Weakness**: System starts executing a test goal every time it boots.
- **Gap**: System should initialize and wait for real external input (or load existing goals from DB).

## 7. Workspace Isolation
- **Current State**: `WorkspaceManager` creates directories and databases but doesn't strictly enforce isolation in the runtime path.
- **Weakness**: Potential for data leakage if runtime components don't explicitly pass and respect `workspace_id`.
- **Gap**: Enforce isolation at the database layer and ensure all components are workspace-aware.
