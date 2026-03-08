# Orion Architecture GAP REPORT

## Existing Systems

The following core runtime and cognitive subsystems have been bootstrapped and are present in the repository:

- **Kernel**: `internal/runtime/kernel.go` - Orchestrates system lifecycle and services.
- **Lifecycle Manager**: `internal/runtime/lifecycle.go` - Handles graceful startup/shutdown.
- **Event Bus**: `internal/runtime/event_bus.go` - Asynchronous, persistent communication with panic recovery.
- **Worker Pool & Scheduler**: `internal/execution/` - Concurrent job processing units.
- **Workspace Manager**: `internal/workspace/` - Handles isolated workspaces with per-workspace SQLite databases.
- **Cognition Engine (OODA-L)**: `internal/cognition/` - Skeleton implementation of the 5-phase cognitive loop.
- **Goal Runtime**: `internal/runtime/goal.go` - Manages persistence and state of user goals.
- **Tool Registry**: `internal/tools/` - Pluggable tool execution framework.
- **SQLite Storage Layer**: `internal/storage/sqlite/` - Support for global and workspace schemas with custom driver support for `sqlite_vec`.
- **Initial Agent Stubs**: Six core agents (Conversation, Planner, Code Indexer, Analysis, Memory Gardener, Pattern Detector) are implemented as stubs and registered in `cmd/orion/main.go`.

---

## Missing Systems

The following functional subsystems are defined in the target architecture but currently exist only as stubs or empty directories:

- **Retrieval Engine**:
    - **Vector Search**: Implementation of `sqlite_vec` integration for semantic retrieval.
    - **Graph Traversal**: Logic to navigate the Zettelkasten-style memory graph.
    - **Hybrid Scoring**: The ranking formula (`score = 0.55 semantic + 0.15 graph + 0.15 temporal + 0.10 usage + 0.05 importance`).
- **Memory Gardener**: Background logic for deduplication, consolidation, and pruning of memory nodes.
- **Code Intelligence**:
    - **Symbol Indexer**: Actual parsing and indexing of Go, Python, and TypeScript source code.
    - **Call Graph Generator**: Construction of cross-file function dependency graphs.
- **Advanced Agent Logic**: The internal reasoning loops for all six agents are currently minimal stubs.
- **Desktop UI**: A Fyne-based interface for agent workspace views, chat, and memory exploration.
- **LLM Prompt Envelope Builder**: Logic to assemble the retrieved context into prompt formats.

---

## Architecture Conflicts

- **Extension Loading**: The current `go-sqlite3` custom driver registration expects `sqlite_vec` to be available in the library path. This might fail on environments without the binary extension, requiring a more robust delivery mechanism for the extension itself.
- **Deterministic vs LLM Decider**: The `Decide` phase currently uses a deterministic stub; a robust interface is needed to toggle between rule-based and LLM-based planning without breaking the pipeline.

---

## Implementation Plan

The following implementation phases are recommended to complete the Orion Cognitive Runtime:

### Phase 1: Core Intelligence Foundation
- Implement `internal/code/parser` using Tree-Sitter for Go, Python, and TypeScript.
- Implement the basic `internal/code/indexer` to populate `code_symbols` and `call_graph_edges` tables.

### Phase 2: Hybrid Retrieval Engine
- Fully integrate `sqlite_vec` into `internal/retrieval/vector`.
- Implement graph traversal logic in `internal/retrieval/graph`.
- Develop the hybrid scoring mechanism in `internal/retrieval/scoring`.

### Phase 3: Memory & Knowledge Management
- Implement the `MemoryGardenerAgent` with background deduplication and consolidation logic.
- Finalize the `Learn` phase of the OODA-L loop to correctly extract patterns and update embeddings.

### Phase 4: Agent Reasoning & LLM Integration
- Upgrade the `PlannerAgent` and `Decide` phase to handle complex plan generation.
- Implement the LLM prompt envelope builder in `internal/cognition/orient.go`.

### Phase 5: Desktop UI
- Initialize the Fyne-based application in `desktop/ui`.
- Create the Workspace View and Chat Interface.

### Phase 6: System Hardening
- Add exhaustive telemetry and logging to the Event Bus and Kernel.
- Implement security boundaries for tool execution.
