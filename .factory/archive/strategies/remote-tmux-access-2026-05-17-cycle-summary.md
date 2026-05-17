---
tags:
  - factory
  - strategy
  - remote-tmux-access
  - cycle-summary
date: 2026-05-17
source: factory-archivist
---

# Cycle Summary: remote-tmux-access — 2026-05-17

## Outcome

**BUILD COMPLETE** — All 8 phases delivered successfully. 100% keep rate (8/8 kept, 0 reverts).

## Final Metrics

- **Total Go code**: 1503 lines
- **Total tests**: 25 passing
- **Phases**: 8/8 complete
- **Keep rate**: 100%
- **Commits**: 9 (871e6dd → cb30c92)

## What Was Built

**rta (Remote tmux Access)** — A single-binary Go CLI that discovers tmux sessions running Claude Code and attaches with one keypress. Designed for accessing Claude Code sessions from an iPhone over SSH.

### Phase Breakdown

| Phase | Feature | Lines | Tests | Commit |
|-------|---------|-------|-------|--------|
| 1 | Project scaffold (Go module, Cobra, Makefile) | scaffold | — | 871e6dd |
| 2 | tmux session/pane parsing | 273 | 11 | 862a9fd |
| 3 | Process tree walking + Claude detection | 204 | 8 | 557bc43 |
| 4 | Cobra subcommands (attach, status, setup) | 230 | — | 1583679, c1559a4 |
| 5 | FindSession + attach edge cases | 96 | 6 | c55148c |
| 6 | Bubble Tea TUI session picker | 309 | — | 9387f32 |
| 7 | Shell profile injection (auto-launch on SSH) | 340 | 11 | c554073 |
| 8 | SSH config checker | 110 | — | b002dfe |

## Key Design Decisions

1. **Go over Python** — Single-binary distribution, zero runtime dependencies
2. **SSH-native, no web UI** — Uses native SSH terminal access via iPhone SSH apps (Terminus, Blink Shell, Prompt 3)
3. **Claude Code detection only** — v1 scoped to finding `claude` binary in tmux pane process trees
4. **Discovery + attachment only** — Never starts, stops, wraps, or sends input to Claude Code
5. **No state on disk** — tmux is sole source of truth
6. **Deferred Charm dependencies** — bubbletea/lipgloss added only at Phase 6, keeping early phases clean

## Architecture

```
cmd/
  root.go        — Cobra root + TUI launch
  attach.go      — Direct attach by name
  status.go      — List sessions (table or --json)
  setup.go       — Setup subcommand group
  setup_ssh.go   — SSH config validation
internal/
  tmux/           — Session/pane parsing, FindSession, AttachSession
  process/        — Process tree walking, Claude detection
  tui/            — Bubble Tea v2 session picker
  profile/        — Shell profile injection, SSH config checker
```

## Strategy Evolution

- **Pre-build research** recommended Python + WebSocket + xterm.js. User pivoted to Go + SSH-native CLI.
- Research findings on competitive landscape and tmux interaction patterns remained valuable despite tech stack pivot.
- 8-phase dependency-ordered plan was created and executed without modification.
- No phase required rework or strategy adjustment — the initial decomposition matched the code's actual dependency graph.

## Patterns Discovered

1. **Research-to-Spec Pivot** — Research technology recommendations are the most likely part to be overridden; emphasize *what* and *why* before *how*.
2. **Implementation Research as Executable Documentation** — Post-pivot research should produce copy-pasteable code samples with exact APIs.
3. **Platform-Specific Error Strings** — CLI tool error messages vary by platform; test error paths on all targets.
4. **Incremental Build With Zero Reverts** — Dependency-ordered phases with one feature per phase and immediate test/lint gates produce reliable incremental progress.

## What's Next (Future Cycles)

Potential v2 features not built in this cycle:
- Multi-agent detection (not just Claude Code)
- Linux `/proc` filesystem support for process tree walking
- Remote session forwarding (SSH tunnel management)
- Session health monitoring / reconnection
- Configuration file support
