---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 4
verdict: KEEP
score_delta: "process tree complete → cobra commands complete"
date: 2026-05-17
source: factory-archivist
---

# Experiment #4: Phase 4 — Cobra Commands (CLI Routing)

## Hypothesis

Wiring Cobra subcommands (`attach`, `status`, `setup`, `setup ssh`) with real implementations using the tmux and process packages provides a functional CLI layer for session discovery and attachment.

## Result

**KEEP** — Phase 4 implemented as specified in the build plan across 2 commits. 230 lines in `cmd/`, all 19 existing tests pass, `go vet` clean. A macOS-specific tmux error handling bugfix was also applied.

## What Changed

### Commit `1583679`: Add Cobra subcommands

Created 4 new files and updated `cmd/root.go` (201 lines added):

**`cmd/attach.go`** (53 lines):
- `rta attach <name>` with `cobra.ExactArgs(1)` validation
- Lists sessions, does substring match via `strings.Contains(session.Name, arg)`
- If exactly one match: attaches via `tmux.AttachSession()` (syscall.Exec)
- If multiple matches: lists them and exits with error
- If no match: lists available sessions and exits with error

**`cmd/status.go`** (111 lines):
- `rta status [--json]` lists all sessions with Claude detection
- Builds process tree once via `process.BuildTree()`, walks for each pane
- Enriched session info: name, hasClaude, workingDir, attached status
- Plain text output: columnar one-line-per-session format
- `--json` flag: marshals to JSON array, prints to stdout

**`cmd/setup.go`** (23 lines):
- `rta setup [--undo]` placeholder — prints message about profile injection (Phase 7)

**`cmd/setup_ssh.go`** (16 lines):
- `rta setup ssh` placeholder — prints message about SSH checking (Phase 8)

**`cmd/root.go`** (27 lines, updated):
- Root command (no subcommand) calls `showStatus(false)` as TUI placeholder
- All subcommands wired via `init()` functions

### Commit `c1559a4`: macOS tmux error handling fix

- `internal/tmux/tmux.go`: `classifyError()` now matches both `"no server running"` (Linux) and `"error connecting to"` (macOS) as `ErrNoServer`
- `cmd/status.go`: handles `ErrNoServer` and `ErrNoSessions` gracefully — prints "No tmux sessions found." (or `[]` with `--json`) instead of returning an error

## Key Decisions

1. **Substring matching for attach** — `strings.Contains(session.Name, arg)` enables `rta attach web` to match `webapp`
2. **Single process tree build** — `BuildTree()` called once in `showStatus`, reused across all panes via `HasDescendant()`
3. **Graceful no-sessions handling** — `status` command returns success (not error) when no tmux sessions exist, with appropriate output format
4. **macOS error string difference** — tmux on macOS says "error connecting to /tmp/tmux-xxx/default" instead of "no server running on /tmp/tmux-xxx/default"; both now classified as `ErrNoServer`

## Exit Criteria Met

- `rta attach <name>` attaches to a real tmux session (substring match works)
- `rta status` lists sessions with Claude detection columns
- `rta status --json` outputs valid JSON
- `go vet ./...` clean
- All 19 tests pass (no new tests added — Phase 4 is primarily wiring)

## Links

- Project: remote-tmux-access
- Commits: `1583679`, `c1559a4`
