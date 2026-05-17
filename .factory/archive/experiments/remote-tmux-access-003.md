---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 3
verdict: KEEP
score_delta: "tmux parsing complete → process tree complete"
date: 2026-05-17
source: factory-archivist
---

# Experiment #3: Phase 3 — Process Tree Walking and Claude Detection

## Hypothesis

An `internal/process/` package with `BuildTree()`, `HasDescendant()`, and `DetectClaude()` enables detection of Claude Code instances running inside tmux panes by walking the process tree from pane PIDs.

## Result

**KEEP** — Phase 3 implemented exactly as specified in the build plan. All 8 tests pass, `go vet` clean. 204 lines added across 2 files.

## What Changed

Created `internal/process/process.go` and `internal/process/process_test.go` in commit `557bc43`:

### Functions
- `BuildTree()` — runs single `ps -ax -o pid,ppid,comm` call, parses into `children map[int][]int` and `names map[int]string`
- `parseTree(out []byte)` — internal parser: skips header, uses `strings.Fields()` for whitespace-delimited parsing, skips malformed lines
- `HasDescendant(rootPID, target, children, names)` — BFS walk from rootPID looking for process with matching name
- `DetectClaude(rootPID)` — convenience wrapper: builds tree + BFS for "claude"

### Design Patterns
- `var runPS` function variable for testability — tests inject mock `ps` output without subprocess calls
- BFS (not DFS) for tree walk — simpler, no stack overflow risk on deep trees
- Single `ps` call for entire tree — avoids N+1 subprocess calls when checking multiple panes

### Tests (8 test cases)
- `TestParseTree` — validates parent→children mapping and pid→name mapping from mock `ps` output
- `TestParseTreeEmpty` — header-only output produces empty maps
- `TestParseTreeMalformed` — non-numeric PIDs skipped, valid entries still parsed
- `TestHasDescendant` — claude found as descendant of sshd (pid 256) and bash (pid 257)
- `TestHasDescendantRootIsTarget` — root pid itself matching target returns true
- `TestHasDescendantMissingPID` — nonexistent PID returns false (no panic)
- `TestBuildTreeWithMock` — end-to-end BuildTree with injected ps output
- `TestDetectClaudeWithMock` — DetectClaude returns true/false correctly with mock

## Key Decisions

1. **`strings.Fields()` over `strings.Split()`** — handles variable whitespace in `ps` output naturally
2. **Function variable `runPS`** — enables mock injection for unit tests without interfaces or test flags
3. **BFS walk** — as planned, using queue slice for simplicity
4. **Match exact base name `"claude"`** — `ps -o comm` returns base command name, direct string equality

## Exit Criteria Met

- `go test ./internal/process/` — all 8 tests pass
- BFS correctly finds "claude" at any depth (direct child, grandchild, root itself)
- Handles empty tree, missing PIDs, malformed lines
- `go vet ./...` clean

## Links

- Project: remote-tmux-access
- Commit: 557bc43
