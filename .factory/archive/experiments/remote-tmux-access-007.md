---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 7
verdict: KEEP
score_delta: "+340 lines, +11 tests"
date: 2026-05-17
source: factory-archivist
---

# Experiment #7: Phase 7 — Shell profile injection

## Hypothesis
Add `internal/profile/` package with Install/Remove functions that idempotently inject an rta auto-launch block into the user's shell profile (~/.zshrc or ~/.bashrc), so rta starts automatically on SSH login.

## Result
**KEEP** — 340 insertions across 3 files, 11 new tests (25 total), all passing, vet clean.

## What Changed
- **`internal/profile/profile.go`** (114 lines): Marker-based block injection. `ProfilePath()` detects shell from $SHELL env, `InstallTo()`/`RemoveFrom()` operate on arbitrary paths for testability. The injected block checks `$SSH_CONNECTION` and `command -v rta` before launching.
- **`internal/profile/profile_test.go`** (213 lines): 11 tests covering install, idempotency, remove, remove-no-marker, non-existent file, file-without-trailing-newline, and round-trip install→remove cycles.
- **`cmd/setup.go`** (+13 lines): Wired `rta setup` to call `profile.Install()` and `rta setup --remove` to call `profile.Remove()`.

## Key Design Choices
- **Marker-based idempotency**: Uses `# rta: remote tmux access (auto-launch)` and `# end rta` markers to fence the block. Repeat installs are no-ops. Removal strips only the marked block.
- **SSH guard**: Block only fires when `$SSH_CONNECTION` is set, so local terminal sessions are unaffected.
- **Testable split**: Public `InstallTo`/`RemoveFrom` accept arbitrary paths, enabling tests without touching real dotfiles.

## Links
- Project: remote-tmux-access
- Commit: `c554073`
