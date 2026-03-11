# Stack Correction & Hardening Gap Report

## Current State Verification

- **Fyne usage:** VERIFIED PRESENT (in `go.mod`).
- **Wails presence:** VERIFIED MISSING (no `wails.json`).
- **Chi presence:** VERIFIED MISSING.
- **Ent presence:** VERIFIED MISSING.
- **Frontend Stack (SvelteKit/Vite):** VERIFIED MISSING (Current is plain Svelte).
- **Tailwind CSS v4 CLI:** VERIFIED PRESENT (in `package.json`).
- **Iconify Tailwind plugin:** VERIFIED PRESENT (in `package.json`).
- **HTMX / htmx-sse:** VERIFIED PRESENT (vendored in `assets/js/`).
- **Assets Directories:** VERIFIED PRESENT (`assets/css`, `assets/js`, `assets/images`, `assets/fonts`).

## Runtime Hardening Status

- **Goal Persistence:** VERIFIED INCORRECT (SQLite table exists but no Ent ORM or robust transactional flow).
- **Job Persistence:** VERIFIED INCORRECT (SQLite table exists but minimal usage).
- **Event Bus Reliability:** VERIFIED INCORRECT (Blocking/unbounded-ish chan used).
- **Supervisor Hardening:** VERIFIED INCORRECT (Missing heartbeats and crash detection).
- **Retrieval Maturity:** VERIFIED INCORRECT (Mocked results in retrieval engine).
- **Bootstrap Correctness:** VERIFIED INCORRECT (Hardcoded demo goals in `main.go`).

## Required Corrections

1. Purge Fyne from `go.mod`.
2. Integrate Wails shell.
3. Implement Chi router in `internal/api/`.
4. Initialize Ent ORM and schema.
5. Upgrade/Align frontend to SvelteKit.
6. Harden Event Bus with bounded buffers and non-blocking drops.
7. Harden Supervisor with heartbeat monitoring.
8. Transition Goal/Job updates to Ent transactions.

## Production Hardening Verification

- **Transactional transitions:** CORRECTED (Implemented in persistence.go).
- **Event Bus reliability:** CORRECTED (Implemented non-blocking bounded buffers).
- **Supervisor monitoring:** CORRECTED (Implemented heartbeat monitoring).
- **ORM Transition:** CORRECTED (Ent schemas defined and generated).
