---
tags:
  - factory
  - project
  - remote-tmux-access
source: factory-archivist
---

# Factory: remote-tmux-access

## Status
- **State**: building (Phase 3 complete, Phase 4 next)
- **Current Score**: process tree + tmux parsing passing (builds, 19 tests pass, vet clean)
- **Experiments Run**: 3
- **Kept**: 3, **Reverted**: 0

## Project Summary

**rta (Remote tmux Access)** — A single-binary Go CLI that lets you access running Claude Code sessions from an iPhone over SSH. You SSH into your Mac, `rta` launches automatically via a shell profile guard, shows which tmux sessions have Claude Code running, and attaches with one keypress.

## Key Design Decisions

1. **Go over Python** — Research recommended Python for MVP speed, but user chose Go for single-binary distribution and zero runtime dependencies. Proven by Upterm's success with similar architecture.
2. **SSH-native, no web UI** — Unlike every competing tool (tmate, sshx, ttyd, Upterm), rta uses native SSH terminal access. No WebSocket, no browser, no xterm.js. The user's iPhone SSH apps (Terminus, Blink Shell, Prompt 3) are the client.
3. **Claude Code detection only** — v1 focuses exclusively on finding the `claude` binary in tmux pane process trees. Multi-agent detection deferred.
4. **Discovery + attachment only** — rta never starts, stops, wraps, or sends input to Claude Code. It's a pure tmux session discovery and attachment layer.
5. **No state on disk** — tmux is the sole source of truth. No database, no config files.
6. **Defer Charm deps** — bubbletea/lipgloss not added to go.mod until Phase 6, avoiding unused dependency warnings during early phases.

## Architecture

- **Language**: Go 1.24, single static binary
- **TUI**: Bubble Tea v2 (Charm ecosystem) — deferred to Phase 6
- **CLI**: Cobra for subcommand routing
- **tmux interaction**: All via `tmux` CLI subprocess calls (list-sessions, list-panes, attach-session)
- **Process detection**: Walk process tree via `ps` (macOS) or `/proc` (Linux) to find `claude` binary

## Build Progress

| Phase | Description | Status | Lines | Tests | Commit |
|-------|-------------|--------|-------|-------|--------|
| 1 | Scaffold + Go module | **DONE** | scaffold | — | `871e6dd` |
| 2 | tmux parsing (`internal/tmux/`) | **DONE** | 273 | 11 | `862a9fd` |
| 3 | Process tree (`internal/process/`) | **DONE** | 204 | 8 | `557bc43` |
| 4 | Cobra commands (`cmd/`) | Pending | — | — | — |
| 5 | syscall.Exec attach e2e | Pending | — | — | — |
| 6 | Bubble Tea TUI (`internal/tui/`) | Pending | — | — | — |
| 7 | Shell profile injection (`internal/profile/`) | Pending | — | — | — |
| 8 | SSH config checker | Pending | — | — | — |

**Totals:** 477 lines, 19 tests, 3/8 phases complete.

### Dependency Graph

```
Phase 1 (scaffold) ✅
    |
    v
Phase 2 (tmux) ✅ ──> Phase 3 (process tree) ✅
                            |
                            v
                      Phase 4 (cobra)
                            |
                            v
                      Phase 5 (attach e2e)
                            |
                            v
                      Phase 6 (TUI)
                            |
                  +---------+---------+
                  v                   v
            Phase 7 (profile)   Phase 8 (SSH)
```

## Research Sources (Implementation)

- [Bubble Tea TUI patterns](sources/bubble-tea-tui-patterns.md) — Custom list, tea.ExecProcess handoff, Lipgloss styling
- [Process tree walking](sources/process-tree-walking.md) — Single ps call, BFS walk, macOS edge cases
- [tmux CLI parsing](sources/tmux-cli-parsing.md) — Format strings, pipe delimiter, error handling
- [Cobra CLI patterns](sources/cobra-cli-patterns.md) — Project layout, root/attach/status/setup commands
- [syscall.Exec attach](sources/syscall-exec-attach.md) — Process replacement vs child spawning
- [Shell profile injection](sources/shell-profile-injection.md) — Marker-based idempotent block install/remove

## Research Sources (Pre-Implementation)

- [Competing tools survey](sources/competing-tools-survey.md) — tmate, sshx, Upterm, ttyd, GoTTY, pyxtermjs
- [Architecture patterns](sources/architecture-patterns.md) — Control mode, PTY wrapper, relay, direct exec
- [Pitfalls from research](sources/pitfalls-from-research.md) — Resize, platform diffs, session discovery

## Ideation Timeline

- **2026-05-17 04:19** — Researcher agent started; completed at 04:24
- **2026-05-17 04:25** — Distiller agent first pass (Python-based spec); completed at 04:25
- **2026-05-17 04:35** — Distiller refinement after user feedback (iOS/SSH focus); completed at 04:36
- **2026-05-17 13:00** — Distiller final pass (Go rewrite, Claude-only, explicit auto-launch); completed at 13:01
- **2026-05-17** — Strategy approved, spec finalized, implementation research completed
- **2026-05-17** — 8-phase build plan created by strategist, CEO-approved
- **2026-05-17** — Phase 1 scaffold built and CEO-approved (PROCEED)
- **2026-05-17** — Phase 2 tmux parsing built (KEEP, 273 lines, 11 tests)
- **2026-05-17** — Phase 3 process tree built (KEEP, 204 lines, 8 tests)

## Recent Experiments

- Experiment #1 — Phase 1 scaffold (**KEEP**, scaffold complete)
- Experiment #2 — Phase 2 tmux parsing (**KEEP**, 273 lines, 11 tests passing)
- Experiment #3 — Phase 3 process tree (**KEEP**, 204 lines, 8 tests passing)
