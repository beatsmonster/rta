---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 2
verdict: KEEP
score_delta: "scaffold → tmux parsing complete"
date: 2026-05-17
source: factory-archivist
---

# Experiment #2: Phase 2 — tmux Session/Pane Parsing

## Hypothesis

A `internal/tmux/` package with `ListSessions()`, `ListPanes()`, and `AttachSession()` provides the foundation for all session discovery and attachment in later phases.

## Result

**KEEP** — Phase 2 implemented exactly as specified in the build plan. All tests pass, `go vet` clean. 273 lines added across 2 files.

## What Changed

Created `internal/tmux/tmux.go` and `internal/tmux/tmux_test.go` in commit `862a9fd`:

### Types
- `Session{Name, Attached, Windows}` — parsed from `tmux list-sessions`
- `Pane{PID, CurrentPath}` — parsed from `tmux list-panes`

### Functions
- `ListSessions()` — runs `tmux list-sessions -F '#{session_name}|#{session_attached}|#{session_windows}'`, parses pipe-delimited output
- `ListPanes(sessionName)` — runs `tmux list-panes -t <session> -F '#{pane_pid}|#{pane_current_path}'`
- `AttachSession(sessionName)` — `syscall.Exec` replaces Go process with `tmux attach-session -t <name>`
- `classifyError()` — distinguishes `ErrNoServer` from `ErrNoSessions` (both exit code 1)

### Error Handling
- Sentinel errors: `ErrNoServer`, `ErrNoSessions` for caller-side error classification
- Stderr inspection on `exec.ExitError` to distinguish tmux failure modes

### Tests (11 test cases)
- `TestParseSessions` — 6 cases: two sessions, single, empty, whitespace-only, malformed line skipped, session name with spaces
- `TestParsePanes` — 5 cases: two panes, path with pipes (SplitN correctness), empty, non-numeric PID skipped, missing fields skipped
- `TestClassifyError_NoServer`, `TestClassifyError_NoSessions`, `TestClassifyError_Other`

## Key Decisions

1. **Pipe `|` delimiter** — as planned, avoiding colon which appears in file paths
2. **`SplitN` with correct field count** — pane paths can contain pipes; `SplitN(line, "|", 2)` preserves them
3. **Graceful skip on malformed lines** — `continue` instead of error, making parser resilient
4. **`syscall.Exec` for attach** — replaces Go process entirely, shell prompt returns on tmux detach

## Exit Criteria Met

- `go test ./internal/tmux/` — all 11 tests pass
- Parse functions handle empty output, missing fields, and "no server" errors
- `go vet ./...` clean

## Links

- Project: remote-tmux-access
- Commit: 862a9fd
