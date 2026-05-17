---
tags:
  - factory
  - strategy
  - remote-tmux-access
date: 2026-05-17
source: factory-archivist
---

# Strategy: remote-tmux-access — 2026-05-17 (Phase 5 Complete)

## Current State

Phase 5 (attach e2e) is complete. The project is 62.5% through the 8-phase build plan (5/8 phases done).

## What Phase 5 Accomplished

Refactored session matching into a reusable `tmux.FindSession` function and hardened all error paths in the attach command. This is the final piece needed before the TUI can call session discovery and attachment as library functions rather than reimplementing the logic.

## Next Phase: Phase 6 — Bubble Tea TUI

Phase 6 is the most complex remaining phase. It will:
- Add `github.com/charmbracelet/bubbletea` and `lipgloss` to `go.mod`
- Create `internal/tui/` with a session list model
- Replace `showStatus()` in root command with `tea.NewProgram()`
- Use `tea.ExecProcess` for attach handoff (discovered in research)
- Display Claude detection status per session

Key risk: Bubble Tea v2 API — research notes cover the patterns, but integration testing will be important.

## Remaining Phases

| Phase | Description | Status | Risk |
|-------|-------------|--------|------|
| 6 | Bubble Tea TUI | Next | Medium — largest phase, new dependency |
| 7 | Shell profile injection | Pending | Low — well-researched pattern |
| 8 | SSH config checker | Pending | Low — read-only checks |

## Metrics

- **Lines of code**: 788 total
- **Tests**: 25 passing
- **Phases**: 5/8 complete (62.5%)
- **Experiments**: 5 run, 5 kept, 0 reverted
