---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 3
verdict: keep
score_delta: "+0.0000"
date: 2026-05-19
source: factory-archivist
---

# Experiment #3: ESC to quit + slog structured logging + 100% test coverage

## Hypothesis
ESC to quit + slog structured logging + 100% test coverage for all packages

## Result
**KEEP** — score delta +0.0000 (precheck scope false positive overridden by CEO)

## What Changed
- **ESC key to quit**: Added `tea.KeyEsc` handler in the Bubble Tea TUI so users can press ESC to exit the session picker
- **slog structured logging**: Added `log/slog` with `--verbose` flag using `PersistentPreRunE` on the Cobra root command; `TextHandler` writes to stderr, Debug level toggled by flag
- **Near-100% test coverage**: Added 1200+ lines of tests across all packages, bringing coverage to 93-100%. Covers `internal/tmux`, `internal/process`, `internal/tui`, `internal/profile`, and `cmd/`
- **Race-safe design**: Unexported mutable function-level variables (`execCommand`, `attachFunc`, etc.) to prevent cross-test contamination; `rootCmd` state reset between tests
- **KeyPressMsg comment**: Added explanatory comment for the custom `KeyPressMsg` type used in TUI testing

## Key Technical Details
- 1200+ lines of new test code
- Function-var injection pattern for testability without mocking frameworks
- `rootCmd` reconstructed per test to avoid shared state across test cases
- PR went through 2 review iterations: first pass flagged shared rootCmd state and exported mutable globals; second pass confirmed fixes

## Links
- Project: remote-tmux-access
- Issue: #5
- PR: #6
- Commits: `2d471ef`, `9fdbf67`
