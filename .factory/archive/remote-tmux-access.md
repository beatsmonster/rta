---
tags:
  - factory
  - project
  - remote-tmux-access
source: factory-archivist
---

# Factory: remote-tmux-access

## Status
- **State**: optimization cycle BLOCKED — score ceiling unreachable for Go projects
- **Current Score**: 0.5011 (threshold: 0.80, max achievable ~0.56)
- **Score Breakdown**: tests=1.0, lint=0.5, type_check=0.5, coverage=0.5, guard_patterns=0.75, config_parser=1.0, capability_surface=0.05, experiment_diversity=0.5, observability=0.0, research_grounding=0.52, factory_effectiveness=0.5
- **Build**: feature-complete CLI — 1503+ lines Go, 25 tests passing, 8/8 phases, 100% keep rate
- **Optimization Experiments**: 2 (both REVERT)
- **Build Experiments**: 8 (all KEEP)
- **Key Constraint**: Factory eval is Python-centric — score ceiling ~0.56 for Go projects. All optimization experiments will be reverted until eval gains Go support or threshold is lowered.

## Current Strategy (2026-05-19)

**BLOCKED** — No viable hypotheses remain within the factory's scoring system.

Two hypotheses were attempted, both reverted:

1. **H2 (eval resilience):** Fix eval/score.py to handle missing Go gracefully — **REVERT** (eval_immutable guard violation, eval/score.py outside declared scope)
2. **H1 (observability):** Add slog structured logging with --verbose flag — **REVERT** (code was CLEAN but score 0.5011 < 0.80 threshold; threshold unreachable for Go)

**Root cause**: The factory eval cannot score Go projects above ~0.56 due to Python-only dimensions (capability_surface, lint, type_check, coverage, observability) collectively weighted at ~0.60.

## Project Summary

**rta (Remote tmux Access)** — A single-binary Go CLI that lets you access running Claude Code sessions from an iPhone over SSH. You SSH into your Mac, `rta` launches automatically via a shell profile guard, shows which tmux sessions have Claude Code running, and attaches with one keypress.

## Key Design Decisions

1. **Go over Python** — Research recommended Python for MVP speed, but user chose Go for single-binary distribution and zero runtime dependencies. Proven by Upterm's success with similar architecture.
2. **SSH-native, no web UI** — Unlike every competing tool (tmate, sshx, ttyd, Upterm), rta uses native SSH terminal access. No WebSocket, no browser, no xterm.js. The user's iPhone SSH apps (Terminus, Blink Shell, Prompt 3) are the client.
3. **Claude Code detection only** — v1 focuses exclusively on finding the `claude` binary in tmux pane process trees. Multi-agent detection deferred.
4. **Discovery + attachment only** — rta never starts, stops, wraps, or sends input to Claude Code. It's a pure tmux session discovery and attachment layer.
5. **No state on disk** — tmux is the sole source of truth. No database, no config files.
6. **Defer Charm deps** — bubbletea/lipgloss not added to go.mod until Phase 6, avoiding unused dependency warnings during early phases. Strategy validated — clean addition in Phase 6.
7. **stdlib slog for logging** — Uses Go's built-in log/slog (no external deps), TextHandler to stderr, --verbose flag for debug output. Added in optimization experiment #2 (reverted due to score threshold, not code quality).

## Architecture

- **Language**: Go 1.24, single static binary
- **TUI**: Bubble Tea v2 + Lipgloss v2 (Charm ecosystem)
- **CLI**: Cobra for subcommand routing
- **Logging**: log/slog with --verbose flag (TextHandler to stderr) — implemented but reverted
- **tmux interaction**: All via `tmux` CLI subprocess calls (list-sessions, list-panes, attach-session)
- **Process detection**: Walk process tree via `ps` (macOS) or `/proc` (Linux) to find `claude` binary
- **Profile injection**: Marker-based idempotent block install/remove in ~/.zshrc or ~/.bashrc
- **SSH validation**: sshd_config parsing to verify remote access compatibility

## Build Progress

| Phase | Description | Status | Lines | Tests | Commit |
|-------|-------------|--------|-------|-------|--------|
| 1 | Scaffold + Go module | **DONE** | scaffold | — | `871e6dd` |
| 2 | tmux parsing (`internal/tmux/`) | **DONE** | 273 | 11 | `862a9fd` |
| 3 | Process tree (`internal/process/`) | **DONE** | 204 | 8 | `557bc43` |
| 4 | Cobra commands (`cmd/`) | **DONE** | 230 | — | `1583679`, `c1559a4` |
| 5 | Attach e2e (FindSession + edge cases) | **DONE** | 96 | 6 | `c55148c` |
| 6 | Bubble Tea TUI (`internal/tui/`) | **DONE** | 309 | — | `9387f32` |
| 7 | Shell profile injection (`internal/profile/`) | **DONE** | 340 | 11 | `c554073` |
| 8 | SSH config checker (`internal/profile/ssh.go`) | **DONE** | 110 | — | `b002dfe` |

**Totals:** 1503+ lines Go, 25 tests, 8/8 phases complete (100%).

### Dependency Graph

```
Phase 1 (scaffold)
    |
    v
Phase 2 (tmux) --> Phase 3 (process tree)
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

## Research Sources (Optimization Cycle — 2026-05-19)

- [Go env constraint](sources/go-env-constraint.md) — Go not installed in eval env, tests/lint/build all fail
- [Factory eval Python bias](sources/factory-eval-python-bias.md) — capability_surface, lint, coverage all Python-only
- [slog observability strategy](sources/slog-observability-strategy.md) — Primary actionable dimension, score formula, recommended approach
- [Backlog items are constraints](sources/backlog-items-are-constraints.md) — Backlog is informational, effective backlog is empty

## Ideation Timeline

- **2026-05-17 04:19** — Researcher agent started; completed at 04:24
- **2026-05-17 04:25** — Distiller agent first pass (Python-based spec); completed at 04:25
- **2026-05-17 04:35** — Distiller refinement after user feedback (iOS/SSH focus); completed at 04:36
- **2026-05-17 13:00** — Distiller final pass (Go rewrite, Claude-only, explicit auto-launch); completed at 13:01
- **2026-05-17** — Strategy approved, spec finalized, implementation research completed
- **2026-05-17** — 8-phase build plan created by strategist, CEO-approved
- **2026-05-17** — Phases 1-8 all built and kept (100% keep rate)

## Cycle Summaries

**Optimization cycle 2026-05-19 COMPLETE.** 2 experiments attempted, both reverted. Score ceiling (~0.56) unreachable for Go projects vs 0.80 threshold. See [full cycle summary](strategies/remote-tmux-access-2026-05-19-cycle-summary.md).

**Build cycle 2026-05-17 COMPLETE.** All 8 phases delivered with 100% keep rate. See [full cycle summary](strategies/remote-tmux-access-2026-05-17-cycle-summary.md).

## Strategy History

- **2026-05-19** — Optimization strategy approved: H2 (eval resilience) then H1 (slog logging). Both reverted. See [strategy snapshot](strategies/remote-tmux-access-2026-05-19.md).
- **2026-05-17** — 8-phase build plan approved and executed. See [build plan](strategies/remote-tmux-access-2026-05-17-build-plan.md).

## Recent Experiments (Optimization Cycle)

- Experiment #1 (opt) — Handle missing Go in eval/score.py (**REVERT**, scope violation: eval/score.py outside declared scope, score 0.426 unchanged)
- Experiment #2 (opt) — Add slog structured logging with --verbose flag (**REVERT**, code was CLEAN but score 0.5011 < 0.80 threshold; systemic — threshold unreachable for Go projects)

## Recent Experiments (Build Cycle)

- Experiment #1 — Phase 1 scaffold (**KEEP**, scaffold complete)
- Experiment #2 — Phase 2 tmux parsing (**KEEP**, 273 lines, 11 tests passing)
- Experiment #3 — Phase 3 process tree (**KEEP**, 204 lines, 8 tests passing)
- Experiment #4 — Phase 4 Cobra commands (**KEEP**, 230 lines, macOS tmux bugfix)
- Experiment #5 — Phase 5 attach e2e (**KEEP**, 96 lines, 6 tests, FindSession reusable)
- Experiment #6 — Phase 6 Bubble Tea TUI (**KEEP**, 309 lines, session picker + auto-attach)
- Experiment #7 — Phase 7 shell profile injection (**KEEP**, 340 lines, 11 tests, marker-based idempotency)
- Experiment #8 — Phase 8 SSH config checker (**KEEP**, 110 lines, sshd_config validation — FINAL PHASE)
