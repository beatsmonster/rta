---
tags:
  - factory
  - strategy
  - remote-tmux-access
date: 2026-05-19
source: factory-archivist
---

# Strategy: remote-tmux-access — 2026-05-19

## CEO Verdict: PROCEED

Both hypotheses approved. No issues found. Priority order: H2 first (unblocks H1 measurement), then H1.

## Current State

- **Composite Score:** 0.426 (threshold: 0.80)
- **Weakest Dimensions:** observability (0.0), capability_surface (0.05, BLOCKED)
- **Critical Constraint:** Go is NOT installed in the eval environment

## Approved Hypotheses

### H2: Make eval/score.py resilient to missing Go toolchain (Priority 1)
- **Category:** FIX
- **Growth dimension:** factory_effectiveness
- **What:** Update eval/score.py so eval_tests() and eval_lint() handle missing Go gracefully (return 0.5 with explanation instead of crashing). This unblocks the observability eval from contributing.
- **Expected impact:** Eval becomes runnable in Go-less environments; observability score can be reported.

### H1: Add slog structured logging across all packages (Priority 2)
- **Category:** EXPLORE
- **Growth dimension:** observability
- **What:** Add log/slog structured logging with --verbose/-v flag in cmd/root.go. Add slog.Debug()/slog.Info() calls across all packages targeting 50%+ function coverage.
- **Expected impact:** observability 0.0 -> ~0.56. Composite score +0.056.

## Key Design Decisions

- Execute H2 before H1 — H2 unblocks H1's eval measurement
- H2 modifies only eval/score.py; H1 modifies cmd/root.go and internal/**/*.go
- Backlog items are constraint notes, not actionable tasks
- Two new backlog items added: trace propagation, capability_surface extension (within cap)

## Anti-patterns Identified

- Don't try to fix capability_surface (factory only counts Python via AST)
- Don't try to fix lint/type_check/coverage scores (no Go integration in factory eval)
- Don't assume Go is available in eval environment
- Don't treat backlog items as tasks

## New Backlog Items

- Add OpenTelemetry trace context propagation (trace_id per CLI invocation)
- Investigate whether factory's capability_surface eval can support Go via project-level override
