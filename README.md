# Orion — Cognitive Runtime for AI Agents

<!-- TODO: Add Logo -->

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org/)
[![Release](https://img.shields.io/badge/release-v0.1--alpha-orange.svg)](https://shields.io/)

---

Orion is a **local cognitive runtime for AI agents** that implements an OODA-L cognition loop, workspace-scoped knowledge orchestration, and autonomous tool execution.

It serves as the **cognitive layer** for agents, managing long-term memory, performing contextual retrieval, and constructing prompt envelopes for LLMs.

---

# Core Philosophy

Orion treats agent memory as a **structured knowledge graph**, not just a collection of chat logs.

By implementing the **OODA-L** loop (Observe, Orient, Decide, Act, Learn), Orion transforms agents from stateless responders into persistent, learning entities.

---

# Key Features

### Local Multi-Agent Runtime

Orion is a desktop application written in Go. It provides a robust, concurrent environment for multiple agents to operate within isolated workspaces.

### Workspace Isolation

Each workspace is a self-contained cognitive environment with its own SQLite database, memory graph, and data artifacts.

### OODA-L Cognition Loop

- **Observe**: Capture user intent and environmental signals.
- **Orient**: Hybrid retrieval across vector space, knowledge graph, and code symbols.
- **Decide**: Deterministic or LLM-based execution planning and tool selection.
- **Act**: Safe execution of tools via a standardized framework.
- **Learn**: Knowledge extraction, pattern detection, and memory consolidation.

### Memory Architecture

Orion uses a **Zettelkasten-style memory graph** where nodes represent facts, insights, and events, linked by semantic and causal relationships.

### Hybrid Retrieval Engine

Retrieval combines:
- **Vector Similarity** (via `sqlite_vec`)
- **Graph Traversal**
- **Code Symbol Search**
- **Temporal & Importance Signals**

---

# Architecture

Orion is built as a highly concurrent Go runtime kernel.

### Components

- **Kernel**: Manages system lifecycle and service orchestration.
- **Event Bus**: Asynchronous, persistent communication between agents.
- **Worker Pool**: High-performance execution units for cognitive tasks.
- **Workspace Manager**: Handles isolated data and runtime environments.
- **Cognition Engine**: Implements the core OODA-L loop phases.
- **Tool Registry**: Pluggable framework for agent capabilities.

### Storage Stack

- **Global Registry**: `data/orion.db` (workspaces, agents, system config).
- **Workspace DB**: `data/workspaces/{id}/workspace.db` (memory, goals, code symbols).
- **Vector Search**: Integrated `sqlite_vec` for high-performance embeddings.

---

# Use Cases

Orion can be used to build:

• persistent AI assistants  
• autonomous agents  
• research copilots  
• long-term conversation systems  
• knowledge bases for LLMs  

Because Orion is **local-first**, it is also suitable for systems requiring:

- privacy
- offline operation
- deterministic behavior

---

# Roadmap

Planned development areas include:

### Cognitive Memory

- semantic timeline queries
- memory graph visualization
- improved memory consolidation

### Agent Infrastructure

- multi-agent coordination
- tool execution framework
- agent workflow orchestration

### Knowledge Systems

- document ingestion
- cross-memory linking
- knowledge evolution

---

# Vision

The long-term goal of Orion is to provide AI systems with **human-like episodic memory**.

Instead of rediscovering context in every prompt, agents can recall prior knowledge naturally.

This transforms AI systems from:

**stateless responders**

into

**persistent cognitive entities**.

---

# License

MIT
