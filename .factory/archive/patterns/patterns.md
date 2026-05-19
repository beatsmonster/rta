---
tags:
  - factory
  - patterns
source: factory-archivist
---

# Cross-Project Patterns

## Research-to-Spec Pivot
Discovered in remote-tmux-access ideation (2026-05-17).
Research recommended a Python + WebSocket + xterm.js web-based architecture. User feedback pivoted the project to Go + SSH-native CLI/TUI with zero web components. The research was still valuable — it identified the key differentiator (bridging to existing tmux sessions) and mapped the competitive landscape — but the technology recommendations were entirely replaced. **Lesson**: Research should emphasize *what* to build and *why* before *how* to build it. Technology recommendations are the most likely part to be overridden by user preferences.

## Implementation Research as Executable Documentation
Discovered in remote-tmux-access implementation research (2026-05-17).
When research produces code samples for a specific tech stack (Go + Bubble Tea + Cobra + tmux), the samples should be directly copy-pasteable — complete imports, correct types, real API calls. The researcher's implementation findings for rta were immediately usable because they included exact function signatures, message types, and edge case handling. **Lesson**: Post-pivot implementation research should produce code-level patterns, not just architectural descriptions. The strategist can then sequence them without re-researching APIs.

## Platform-Specific Error Strings Require Real-Device Testing
Discovered in remote-tmux-access Phase 4 (experiment #4, 2026-05-17).
tmux on macOS reports "error connecting to /tmp/tmux-xxx/default" when no server is running, while Linux reports "no server running on /tmp/tmux-xxx/default". The research source notes warned about platform differences, but the exact error string divergence was only caught during real-device testing. **Lesson**: When wrapping CLI tools that produce human-readable error messages, test error paths on all target platforms. Error string matching is fragile — check for multiple known variants or use exit codes where possible.

## Incremental 8-Phase Build With Zero Reverts
Discovered in remote-tmux-access full build (2026-05-17).
An 8-phase incremental build plan achieved a 100% keep rate (8/8 phases kept, 0 reverts). Key factors: (1) each phase had a clear dependency on prior phases, preventing integration surprises; (2) deferred optional dependencies (Charm libs) until the phase that needed them; (3) platform-specific bugs were caught and fixed within the phase that introduced them rather than accumulating. **Lesson**: For greenfield CLI projects, a dependency-ordered phase plan with one feature per phase and immediate test/lint gates produces reliable incremental progress. The keep rate tracks build quality — 100% means the strategy's phase decomposition matched the code's actual dependency graph.

## Factory Eval Python Bias Creates Hard Score Ceiling for Non-Python Projects
Discovered in remote-tmux-access eval analysis (2026-05-19).
The factory's eval system (`hygiene.py` + `growth.py`) is Python-centric: `capability_surface` only counts Python files via AST, lint only checks ruff/eslint/clippy (no Go), coverage only runs pytest. For a Go project with 18 source files and 25 tests, these dimensions all return 0.5 ("not detected") or near-zero. Combined weight of blocked dimensions: ~0.40. Maximum achievable score without factory changes: ~0.56 (threshold: 0.80). **Lesson**: When the factory eval doesn't support a project's language, focus hypotheses exclusively on dimensions that scan source files directly (observability, research_grounding) rather than those requiring compilation or test execution. The score ceiling should be documented upfront so the strategist doesn't waste cycles on structurally impossible targets.

## Missing Toolchain in Eval Environment Masks True Quality
Discovered in remote-tmux-access eval analysis (2026-05-19).
Go was not installed in the eval environment, causing `go test ./...` and `go vet ./...` to fail silently. The factory reports tests as "not detected" (0.5) despite 25 passing tests across 3 test files. **Lesson**: When eval scores show "not detected" for dimensions that clearly have corresponding project artifacts, investigate whether the eval environment has the required toolchain installed. The eval/score.py should handle missing tools gracefully (return partial scores with diagnostic messages) rather than silently scoring 0.5.

## Eval Infrastructure Is Immutable — Don't Modify eval/score.py
Discovered in remote-tmux-access experiment #1 (optimization, 2026-05-19).
The experiment attempted to fix eval/score.py to handle missing Go toolchain gracefully. The change was technically correct (+16 lines, two `except FileNotFoundError` handlers) and passed CEO code review as CLEAN on all 7 checks. However, it was **reverted** because eval/score.py is outside the project's declared scope (`main.go`, `cmd/**/*.go`, `internal/**/*.go`, `Makefile`, `CLAUDE.md`). The eval_immutable guard exists to prevent experiments from gaming their own scoring system. **Lesson**: Never modify eval infrastructure as part of an experiment, even when the change is well-intentioned. If the eval has a bug, it must be fixed through a separate process outside the factory experiment loop. Strategists should not generate hypotheses that target eval files.

## stdlib slog Is Zero-Cost Observability for Go Projects
Discovered in remote-tmux-access experiment #2 (optimization, 2026-05-19).
Adding `log/slog` structured logging across 8 files required only +47 lines with zero new dependencies (stdlib only). The pattern: `PersistentPreRunE` on Cobra root command initializes a `slog.TextHandler` writing to stderr; `--verbose` flag toggles between Info and Debug levels. Debug calls are no-ops at Info level, so there's no performance cost in normal usage. CEO review was CLEAN on all 7 checks. **Lesson**: For Go CLI projects, slog is the lowest-friction path to observability. Initialize once in the root command, sprinkle `slog.Debug` at key decision points (command execution, external tool calls, parsing results), and the observability dimension improves with no behavioral changes to test.

## Clean Code Can Be Reverted by Systemic Eval Limitations
Discovered in remote-tmux-access experiment #2 final verdict (optimization, 2026-05-19).
Experiment #2 (slog logging) passed CEO code review as CLEAN on all 7 checks, was merged as PR #4, and improved the score from 0.426 to 0.5011. However, it was still reverted because 0.5011 < 0.80 threshold. The threshold is structurally unreachable — Python-only eval dimensions (capability_surface, lint, type_check, coverage) carry ~0.40 combined weight and cannot score above 0.5 for Go projects. The observability dimension scored 0.0 even after adding slog, suggesting it checks for Python logging patterns. **Lesson**: When the eval system's score ceiling is below the keep threshold, the optimization loop becomes a dead end. The factory should detect this condition (multiple reverts of CLEAN code) and escalate rather than continuing to generate hypotheses. Strategists should calculate the theoretical max score *before* proposing optimization experiments, and if max < threshold, recommend an eval system update instead.
