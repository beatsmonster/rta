---
tags:
  - factory
  - project
  - remote-tmux-access
source: factory-archivist
---

# Factory: remote-tmux-access

## Status
- **State**: Feature-complete — all user-requested features delivered
- **Current Score**: 0.94 (Go-centric eval)
- **Score Breakdown**: tests=1.0, coverage=0.95, lint=1.0, observability=0.45
- **Build**: Feature-complete CLI — 2200+ lines Go, 30+ tests passing, 95% avg coverage
- **Total Experiments**: 13 (build cycle) + 5 (optimization/feature cycle) = 18
- **Kept**: 11 (8 build + 3 feature), **Reverted**: 2

## Current Cycle (2026-05-19) — COMPLETE

5 experiments run: 2 reverted (eval bias), 3 kept (feature delivery). All 3 user-requested features shipped in PR #6:

1. **ESC to quit + slog + 100% test coverage** (exp 3, KEEP)
2. **Recap infrastructure + TUI display** (exp 4, KEEP)
3. **/enable-rta skill + setup subcommand** (exp 5, KEEP)

## Project Summary

**rta (Remote tmux Access)** — A single-binary Go CLI that lets you access running Claude Code sessions from an iPhone over SSH. You SSH into your Mac, `rta` launches automatically via a shell profile guard, shows which tmux sessions have Claude Code running, and attaches with one keypress.

## Key Design Decisions

1. **Go over Python** — Single-binary distribution, zero runtime dependencies
2. **SSH-native, no web UI** — Native terminal via iPhone SSH apps (Terminus, Blink Shell, Prompt 3)
3. **Claude Code detection only** — v1 focuses on finding `claude` binary in tmux pane process trees
4. **Discovery + attachment only** — Pure tmux session discovery and attachment layer
5. **No state on disk** — tmux is the sole source of truth (except recaps in ~/.rta/recaps/)
6. **stdlib slog for logging** — log/slog with --verbose flag, TextHandler to stderr
7. **go:embed for skill bundling** — Skill file compiled into binary, installed via `rta setup skill`

## Architecture

- **Language**: Go 1.24, single static binary
- **TUI**: Bubble Tea v2 + Lipgloss v2 (Charm ecosystem)
- **CLI**: Cobra for subcommand routing
- **Logging**: log/slog with --verbose flag (TextHandler to stderr)
- **tmux interaction**: All via `tmux` CLI subprocess calls
- **Process detection**: Walk process tree via `ps` (macOS) or `/proc` (Linux)
- **Profile injection**: Marker-based idempotent block install/remove in ~/.zshrc or ~/.bashrc
- **SSH validation**: sshd_config parsing for remote access compatibility
- **Recap**: JSON files in ~/.rta/recaps/ for session annotations
- **Skill**: go:embed SKILL.md installed to ~/.claude/skills/

## Build Progress

| Phase | Description | Status | Lines | Tests |
|-------|-------------|--------|-------|-------|
| 1 | Scaffold + Go module | **DONE** | scaffold | — |
| 2 | tmux parsing | **DONE** | 273 | 11 |
| 3 | Process tree | **DONE** | 204 | 8 |
| 4 | Cobra commands | **DONE** | 230 | — |
| 5 | Attach e2e | **DONE** | 96 | 6 |
| 6 | Bubble Tea TUI | **DONE** | 309 | — |
| 7 | Shell profile injection | **DONE** | 340 | 11 |
| 8 | SSH config checker | **DONE** | 110 | — |
| 9 | ESC + slog + test coverage | **DONE** | 1200+ | 25+ |
| 10 | Recap infrastructure | **DONE** | 719 | 10+ |
| 11 | /enable-rta skill | **DONE** | — | — |

## Research Sources

- [Bubble Tea TUI patterns](sources/bubble-tea-tui-patterns.md)
- [Process tree walking](sources/process-tree-walking.md)
- [tmux CLI parsing](sources/tmux-cli-parsing.md)
- [Cobra CLI patterns](sources/cobra-cli-patterns.md)
- [syscall.Exec attach](sources/syscall-exec-attach.md)
- [Shell profile injection](sources/shell-profile-injection.md)
- [Competing tools survey](sources/competing-tools-survey.md)
- [Architecture patterns](sources/architecture-patterns.md)
- [Pitfalls from research](sources/pitfalls-from-research.md)
- [Go env constraint](sources/go-env-constraint.md)
- [Factory eval Python bias](sources/factory-eval-python-bias.md)
- [slog observability strategy](sources/slog-observability-strategy.md)
- [Backlog items are constraints](sources/backlog-items-are-constraints.md)

## Cycle Summaries

**Feature cycle 2026-05-19 COMPLETE.** 5 experiments (3 kept, 2 reverted). All user-requested features delivered. See [final cycle summary](strategies/remote-tmux-access-2026-05-19-final-cycle-summary.md).

**Optimization cycle 2026-05-19 (early).** 2 experiments attempted, both reverted due to Python-centric eval bias. See [optimization cycle summary](strategies/remote-tmux-access-2026-05-19-cycle-summary.md).

**Build cycle 2026-05-17 COMPLETE.** 8 phases, 100% keep rate. See [build cycle summary](strategies/remote-tmux-access-2026-05-17-cycle-summary.md).

## Experiment History

### Feature Cycle (2026-05-19, post-eval-reconfiguration)
- Experiment #3 — ESC + slog + test coverage (**KEEP**, PR #6)
- Experiment #4 — Recap infrastructure (**KEEP**, PR #6)
- Experiment #5 — /enable-rta skill (**KEEP**, PR #6)

### Optimization Cycle (2026-05-19, pre-eval-reconfiguration)
- Experiment #1 — eval/score.py resilience (**REVERT**, scope violation)
- Experiment #2 — slog structured logging (**REVERT**, threshold unreachable)

### Build Cycle (2026-05-17)
- Experiments #1-8 — Phases 1-8 (**ALL KEEP**, 100% keep rate)
