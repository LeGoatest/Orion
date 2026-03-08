# GAP REPORT - Orion Runtime

## Existing Systems
- Kernel runtime (internal/runtime)
- Workspace manager (internal/runtime/workspace)
- Event bus (internal/types)
- Scheduler (internal/execution/scheduler)
- Worker pool (internal/execution/worker)
- Goal runtime (internal/runtime/goal)
- Memory manager/graph/gardener (internal/memory)
- Retrieval framework (internal/retrieval)
- Tool registry (internal/tools)
- Code parser/indexer scaffolding (internal/code)

## Missing/Partial Systems
- **Pattern Engine**: No logic to detect or reuse successful goal execution patterns.
- **Retrieval Graph Expansion**: Retrieval currently relies on vector search/hybrid scoring but lacks explicit link-based expansion.
- **Symbolic Index**: No deterministic symbol lookup before vector retrieval.
- **Integration**: Cognition loop doesn't yet bypass reasoning for pattern matches or leverage symbolic lookup first.

## Conflict/Redesign Notes
- No architectural redesign requested.
- Adhere to per-workspace database isolation.
- Follow existing OODA-L loop structure.
