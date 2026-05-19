## Strategy — 2026-05-19

### Design Space
| Dimension | Score | Notes |
|---|---|---|
| Features | 5 | All 8 build phases complete: tmux parsing, process tree, CLI, TUI, profile, SSH checker |
| Bug fixes | 2 | One fix merged (macOS tmux error classification), no open issues |
| Instrumentation | 0 | Zero log statements in non-test Go files; observability eval = 0.0 |
| Flow changes | 1 | Core architecture is stable, no refactors attempted |
| New agents | 0 | N/A — not an agent project |
| Prompt engineering | 0 | N/A — not an AI agent project |
| Eval improvements | 1 | eval/score.py exists but crashes when Go is unavailable |
| Knowledge management | 1 | Archive exists but is empty in this worktree |
| Infrastructure | 2 | Makefile exists, factory.md configured |
| Operational execution | 0 | No end-to-end runs recorded |
| Self-evolution | 0 | No factory self-improvement experiments |

**Underserved:** Instrumentation, Eval improvements, Operational execution

### Observations
- Current composite score: **0.426** (threshold: 0.80)
- Weakest eval dimension: **observability** (0.0, weight 0.10)
- Next weakest: **capability_surface** (0.05, weight 0.14) — BLOCKED (factory counts Python only)
- Last 3 experiments: none recorded
- Pattern: Project is feature-complete (all 8 build phases merged) but has zero observability infrastructure and the eval/score.py crashes when Go is not installed. The factory eval also can't detect Go tests/lint/type_check/coverage, so those dimensions are stuck at 0.5.
- **Critical constraint**: Go is NOT installed in the eval environment. `go test`, `go vet`, `go build` all fail with "command not found". Hypotheses must focus on improvements scorable via source code scanning.
- The backlog items are constraint notes (prerequisites, not tasks) — the effective backlog is empty.
- The eval/score.py runs `go test` and `go vet` which crash in this environment. The observability eval function scans source files directly (Python, no Go binary needed) — this is the only project-level eval dimension that can actually improve.
- `guard_patterns` scores 0.75 — 2 pattern tests fail expecting `cmd/example.py` and `internal/example.py` to match guards. These are factory-side test artifacts, not actionable.

### Hypotheses

#### H1: Add slog structured logging across all packages
- **Category:** EXPLORE
- **Growth dimension:** observability
- **New:** yes
- **What:** Add `log/slog` structured logging to rta. Specifically: (1) Add a `--verbose`/`-v` persistent flag to `cmd/root.go` that toggles between `slog.LevelInfo` and `slog.LevelDebug`, initialized in `PersistentPreRunE`. (2) Add `slog.Debug()` and `slog.Info()` calls to key functions across all packages: `tmux.ListSessions`, `tmux.ListPanes`, `process.BuildTree`, `process.HasDescendant`, `profile.Install`, `profile.Remove`, `ssh.CheckSSHConfig`, `tui.refreshSessions`, `cmd.showStatus`, `cmd.attachToSession`. Target: 50%+ function coverage with log statements. Logs go to stderr to keep stdout clean.
- **Why:** Observability is 0.0 — the project has zero logging. The eval/score.py observability function scans Go source files directly for `slog.\w+\(` patterns without needing Go installed. This is the highest-impact dimension we can actually improve. Research confirms `log/slog` (stdlib since Go 1.21) is the standard approach for Go CLI tools.
- **Expected impact:** observability 0.0 → ~0.56 (coverage=0.6, structured=yes, tracing=no, density=0.5). Composite score +0.056.
- **Priority:** high

#### H2: Make eval/score.py resilient to missing Go toolchain
- **Category:** FIX
- **Growth dimension:** factory_effectiveness
- **New:** yes
- **What:** Update `eval/score.py` so the `eval_tests()` and `eval_lint()` functions handle the case where `go` is not found. When `subprocess.run(['go', ...])` raises `FileNotFoundError`, return score 0.5 with details "Go toolchain not available" instead of crashing. This makes the eval script runnable in environments without Go installed, allowing the observability eval to contribute its score.
- **Why:** The CEO's review explicitly calls for "fixing the eval/score.py to handle missing Go gracefully (return partial scores rather than crashing)." Currently if `go` is not installed, the eval script crashes entirely with an unhandled exception, which means even the observability eval (which needs no Go) never gets to report its score. The eval script has 3 functions: `eval_tests`, `eval_lint`, `eval_observability`. The first two call `go` commands; the third scans source files. If the first function crashes, the whole script fails.
- **Expected impact:** factory_effectiveness improvement — eval becomes runnable in Go-less environments. Tests 0.5 stays 0.5, lint 0.5 stays 0.5, but observability score can now be reported. This unblocks H1's impact from being measured.
- **Priority:** high

### Anti-patterns to Avoid
- **Don't try to fix capability_surface**: The factory's growth.py only counts Python files via AST parsing. All Go code is invisible. No amount of rta changes will fix this.
- **Don't try to fix lint/type_check/coverage scores**: The factory eval has no Go linter/type checker/coverage integration. These are stuck at 0.5.
- **Don't assume Go is available**: The eval environment lacks Go. Any hypothesis requiring `go build`, `go test`, or `go vet` will fail at eval time.
- **Don't treat backlog items as tasks**: The current backlog items are constraint notes ("tmux is assumed to be installed"), not actionable work items.

### New Backlog Items
- Add OpenTelemetry trace context propagation (trace_id per CLI invocation) to push observability score higher via the `has_trace` component
- Investigate whether the factory's capability_surface eval can be extended to support Go via a project-level override in eval/score.py
