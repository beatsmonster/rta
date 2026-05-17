---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 6
verdict: KEEP
score_delta: "+309 lines, TUI functional"
date: 2026-05-17
source: factory-archivist
---

# Experiment #6: Phase 6 — Bubble Tea TUI session picker

## Hypothesis
Add an interactive TUI using Bubble Tea v2 and Lipgloss v2 that lists tmux sessions split into Claude Code / Other sections, supports cursor navigation, Enter to attach, r to refresh, and q to quit. Root command auto-attaches when exactly one Claude session exists.

## Result
**KEEP** — Phase 6 complete. TUI fully functional with session picker, auto-attach logic, and Charm ecosystem dependencies added.

## What Changed
- `internal/tui/model.go` (159 lines): Full Bubble Tea model with Init/Update/View, cursor navigation (j/k/arrows), section-aware rendering (Claude Code sessions vs Other sessions), Enter to attach via `tea.ExecProcess`, r to refresh session list, q to quit
- `internal/tui/commands.go` (51 lines): Tea commands for loading sessions (`loadSessions`) and attaching via tmux (`attachSession`) using `tea.ExecProcess` for clean process handoff
- `internal/tui/doc.go` removed (placeholder replaced by real code)
- `cmd/root.go` (+48 lines): Root command now runs TUI by default; auto-attaches when exactly one Claude session exists (skips TUI entirely); integrates `FindSession` for discovery
- `go.mod` / `go.sum`: Added `bubbletea v2`, `lipgloss v2`, and transitive Charm dependencies (as planned — deferred from earlier phases to avoid unused dep warnings)

## Stats
- **Lines added**: 309 (net)
- **TUI package**: 210 lines (commands.go + model.go)
- **Total project**: ~1057 lines Go
- **Tests**: 14 passing (8 process + 6 tmux), vet clean
- **TUI tests**: None yet (TUI is interactive, hard to unit test without mocking tea.Model)

## Key Design Choices
1. **tea.ExecProcess for attach**: Clean process replacement — Bubble Tea suspends, tmux takes over the terminal, Bubble Tea resumes on detach
2. **Auto-attach shortcut**: When exactly one Claude session exists, skip the TUI entirely and attach directly — optimized for the common iPhone SSH case
3. **Section-based rendering**: Sessions grouped into "Claude Code" and "Other" sections, with cursor awareness of section boundaries

## Links
- Project: remote-tmux-access
- Commit: `9387f32`
