---
tags:
  - factory
  - strategy
  - remote-tmux-access
date: 2026-05-17
source: factory-archivist
---

# Strategy: remote-tmux-access — 2026-05-17 (Phase 7 Complete)

## Current State
7 of 8 phases complete (87.5%). Only Phase 8 (SSH config checker) remains.

## Phase 7 Outcome
Shell profile injection implemented cleanly. The `internal/profile/` package adds marker-fenced Install/Remove functions that idempotently manage a 4-line auto-launch block in ~/.zshrc or ~/.bashrc. The block only fires over SSH ($SSH_CONNECTION guard), ensuring local terminals are unaffected. 11 tests cover all edge cases including idempotency, missing files, and round-trips.

## Remaining Work
- **Phase 8: SSH config checker** — Validate that the user's SSH server config permits the rta workflow (e.g., AllowTcpForwarding, TERM passthrough). This is the final phase.

## Running Totals
- **Lines of Go**: ~1397 (1057 + 340)
- **Tests**: 25 (14 + 11)
- **Phases complete**: 7/8
- **Kept**: 7, **Reverted**: 0
