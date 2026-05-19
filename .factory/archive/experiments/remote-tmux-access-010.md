---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 2
verdict: KEEP
score_delta: pending
date: 2026-05-19
source: factory-archivist
---

# Experiment #2 (Optimization Cycle): Add slog structured logging with --verbose flag

## Hypothesis
Add log/slog structured logging across all packages with a --verbose flag to improve observability score from 0.0 to ~0.56.

## Result
**KEEP (CLEAN)** — CEO verdict CLEAN with all checklist items passing. PR #4 merged.

## What Changed
- `cmd/root.go`: Added `--verbose`/`-v` persistent flag on root command; `PersistentPreRunE` initializes slog with `TextHandler` writing to stderr, level toggled between Info and Debug based on flag
- `cmd/attach.go`: Added `slog.Debug` call when attaching to a session
- `cmd/status.go`: Added `slog.Debug` calls for session listing and JSON output paths
- `internal/tmux/tmux.go`: Added `slog.Debug` for tmux command execution and session parsing (14 lines changed)
- `internal/process/process.go`: Added `slog.Debug` for process tree building and Claude detection (10 lines changed)
- `internal/profile/profile.go`: Added `slog.Debug` for profile injection operations
- `internal/profile/ssh.go`: Added `slog.Debug` for SSH config checking
- `internal/tui/commands.go`: Added `slog.Debug` for TUI command execution

**Totals:** 8 files, +47 lines, -4 lines

## CEO Review Highlights
- Correctness: PASS — slog import/usage correct, TextHandler to stderr, PersistentPreRunE initialization
- Security: PASS — no secrets, logging to stderr only
- Edge cases: PASS — defaults to Info level, Debug calls are no-ops in normal mode
- Missing tests: PASS — logging is infrastructure, doesn't alter control flow
- Style: PASS — snake_case attribute keys, consistent slog.Debug usage
- Scope: PASS — all changes within declared scope (cmd/*.go, internal/**/*.go)
- Guardrails: PASS — no eval/score.py or .factory/ modifications

## Design Notes
- Uses Go stdlib `log/slog` — zero new dependencies
- All logs go to stderr to keep stdout clean for structured output (status --json)
- `PersistentPreRunE` on root command ensures logger is initialized before any subcommand runs
- Debug-level logging is silent by default; only enabled with `--verbose`
- snake_case attribute keys throughout for consistency

## Links
- Project: remote-tmux-access
- Issue: #3
- PR: #4
- Strategy: H1 from [2026-05-19 strategy](../strategies/remote-tmux-access-2026-05-19.md)
- Prior experiment: [Experiment #1](remote-tmux-access-009.md) (REVERT — eval scope violation)
