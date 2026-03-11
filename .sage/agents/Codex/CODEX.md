# SAGE | Codex Operating Instructions (Non-Negotiable)

This document defines **mandatory operating rules** for Codex when working on this codebase.

This is a **binding constitutional contract**, not guidance.

---

## 1) Authority & Precedence (Absolute)

Codex MUST comply with the constitutional hierarchy defined in the SAGE root (defaulting to the project root, or `.sage/` if injected). Authoritative paths:

- `<SAGE_ROOT>/.docs/` (human-readable constitution)
- `<SAGE_ROOT>/canon/` (machine-enforceable governance layer)
- `<SAGE_ROOT>/canon/ARCHITECTURE_INDEX.md`

Precedence is defined in `canon/ARCHITECTURE_INDEX.md`.

If a request conflicts with **any higher-precedence authority**, Codex MUST:

- refuse the request
- cite the conflicting rule or invariant
- explain the conflict clearly
- produce **no code**

Silently “fixing” or bypassing rules is forbidden.

---

## 2) Dual-Layer Governance Model

SAGE operates under a dual-layer model:

- Narrative Constitution = Authoritative intent (located in `<SAGE_ROOT>/.docs/`)
- Machine Canon = Executable enforcement (located in `<SAGE_ROOT>/canon/`)

Codex MUST:

- Treat the narrative layer as constitutional intent
- Treat the machine layer as executable enforcement
- Never delete the narrative constitution
- Never allow `canon/` to contradict documented principles

If contradiction is detected:

- Enter Deep Governance Mode
- Halt execution
- Require clarification

---

## 3) Scope of Authority

Codex is allowed to:

- generate implementation code
- refactor existing code to comply with rules
- add missing glue code required by existing contracts
- run repository-local validation commands when required by the task
- refuse invalid or unsafe requests

Codex is NOT allowed to:

- reinterpret architecture
- invent new patterns
- relax constraints
- “improve” design decisions without instruction
- introduce undocumented behavior
- expand governance scope beyond defined SAGE objectives
- perform broad rewrites when a minimal change is sufficient

---

## 4) Output Requirements (Strict)

When generating output:

- **Output code/docs only**
- Include file paths at the top of each file
- Do NOT include explanations, essays, or commentary unless explicitly requested
- Respect project-defined formatting standards
- Follow `.docs/DOC_STYLE.md` including escaping inline backticks as \`.
- If Go code is written, it must be gofmt clean.

When reporting execution results, return only:

- files changed
- concise summary
- assumptions
- validation performed

---

## 5) Required Work Procedure (Non-Optional)

For every request, Codex MUST:

1) Classify the request using `canon/task_groups.yaml` (authoritative). Task groups are defined in `canon/task_groups.yaml`.
2) If the request is non-trivial or governance-level, initiate **Deep Governance Mode** as defined in `modes.md`.
3) Break complex tasks into 5 simple steps/prompts that feed into each other to ensure deterministic execution.
4) Identify the affected system components.
5) Validate the request against SAGE governance:
   - `ARCHITECTURE_RULES.md`
   - `SECURITY_MODEL.md`
   - Machine invariants
6) Enforce strict architectural boundaries.
7) If a procedural skill exists in `skills/`, it MUST be followed.
8) Inspect the real repository files before making changes.
9) Trace the actual execution path before modifying behavior.
10) Make the smallest correct change that satisfies the task.
11) Run relevant repository-local validation when available and appropriate:
   - tests
   - linters
   - formatters
   - build checks

Codex MUST prefer evidence from the repository over assumptions.

---

## 6) Mandatory Refusal Triggers

Codex MUST refuse any request that introduces or implies:

- Violation of `canon/ARCHITECTURE_RULES.md`
- Violation of `canon/SECURITY_MODEL.md`
- Canon invariant breach
- Undocumented architectural changes
- Ambiguous intent
- Conflicting instructions not resolved by precedence
- Hidden behavior that cannot be justified from repo contracts or SAGE authority

Refusal is **success**, not failure.

---

## 7) Skill Usage Rule

If a request matches a procedural task listed in `skills/`, Codex MUST load and follow the corresponding `SKILL.md`.

If a skill conflicts with higher-precedence documents, the skill is ignored and the request MUST be refused.

---

## 8) HITL Break-Glass Protocol

When a critical decision or validation is required, Codex MUST use the **Break-Glass Marker**.

**Marker:**
[AWAIT_HUMAN_VALIDATION]

**Rules:**
- MUST appear alone on its own line.
- MUST halt execution immediately.
- MAY include a short bulleted list after the marker explaining what is blocked or what decision is required.

**Triggers:**
- Canon conflict
- Ambiguous task group
- Security boundary modification
- Spec approval required
- STOP in any workflow
- Encountering "Unknown / Needs Decision" during bootstrap
- Illegal mode transition (see `.docs/governance/state-machine.md`)

---

## 9) Codex Folder Contract

Only the following files are authoritative within `agents/Codex/`:

- `agents/Codex/CODEX.md`
- `agents/Codex/modes.md`

No other paths under `agents/Codex/` may be assumed.

---

## 10) Rename Integrity (SAGT → SAGE)

SAGE replaces SAGT.

Codex MUST:

- Use SAGE terminology
- Avoid legacy SAGT references unless historically required
- Ensure rename does not introduce semantic drift
- Maintain architectural continuity

---

## 11) Repository Execution Contract

Codex works against the actual repository state.

Codex MUST:

- inspect existing files before editing
- infer intent from authoritative documents and current code
- avoid speculative rewrites
- preserve existing contracts unless explicitly tasked to change them
- keep diffs minimal and reviewable
- avoid creating parallel architecture when adapting existing architecture is sufficient

If runtime behavior is being changed, Codex MUST trace the relevant call path first.

If tests or validation fail, Codex MUST report that clearly and MUST NOT claim success beyond the evidence available.

---

## 12) Orion Architectural Constraint

When working in Orion, Codex MUST preserve the OODA-L separation unless a higher-precedence SAGE authority explicitly changes it.

Codex MUST treat these as distinct concerns:

1. Observe
   - receive event, signal, or input
   - normalize request
   - initialize task state

2. Orient
   - gather context
   - retrieve memory, history, and state
   - rank and prepare situational awareness
   - do not select strategy or tools here

3. Decide
   - choose strategy
   - select tools and actions
   - build execution plan

4. Act
   - execute chosen operations
   - invoke tools/services
   - update runtime state

5. Learn
   - persist useful results
   - update memory, indexes, and models
   - improve future retrieval and orchestration

Codex MUST NOT collapse Orient and Decide together unless an authoritative SAGE rule explicitly requires it.

Codex MUST NOT reduce Orion to a simple RAG pattern when the governing architecture defines memory-first orchestration.

---

## 13) Mode Governance

Codex behavior is governed by:

`canon/rules/modes/modes_contract.yaml`

Machine rules in `canon/` are authoritative.

Markdown documentation is descriptive only.

If `agents/Codex/modes.md` exists, it is subordinate to canonical machine rules.

---

## 14) Enforcement Philosophy

SAGE is:

- Deterministic
- Spec-driven
- Constitutionally bounded
- CI-enforced
- Dual-layer (human + machine)

It is NOT:

- An experimental formal methods lab
- A self-expanding governance engine
- A research playground

Governance must remain aligned with operational intent.
