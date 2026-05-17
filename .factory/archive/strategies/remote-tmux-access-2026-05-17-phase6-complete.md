---
tags:
  - factory
  - strategy
  - remote-tmux-access
date: 2026-05-17
source: factory-archivist
---

# Strategy: remote-tmux-access — 2026-05-17 (Phase 6 Complete)

## Current State
Phase 6 (Bubble Tea TUI) is complete. The project now has a fully interactive session picker with auto-attach, built on the Charm ecosystem (bubbletea v2, lipgloss v2). 6 of 8 phases done (75%).

## What's Left

### Phase 7: Shell Profile Injection (`internal/profile/`)
- Idempotent marker-based block install/remove in `.bashrc`/`.zshrc`
- Auto-launch `rta` on SSH login (detect SSH via `$SSH_CONNECTION`)
- Research already done: see `sources/shell-profile-injection.md`

### Phase 8: SSH Config Checker
- Verify `sshd_config` allows the user's setup
- Check for common pitfalls (PasswordAuthentication, AllowTcpForwarding, etc.)
- Lighter phase — may be combined with Phase 7

## Build Velocity
- 6 phases completed in a single day
- ~1057 lines of Go, 14 tests passing
- Clean architecture: each phase builds on the last without rework
- Charm dependency deferral strategy worked perfectly — added in Phase 6 exactly when needed

## Risk Assessment
- **Low risk**: Phases 7-8 are well-researched and relatively straightforward
- **Testing gap**: TUI has no unit tests (interactive components are hard to test); process and tmux packages are well-tested
- **Platform risk**: Shell profile injection needs to handle both bash and zsh on macOS; process tree walking already proven on macOS
