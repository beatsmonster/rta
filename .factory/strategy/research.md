# Research Report

## Project Summary

rta (Remote tmux Access) is a single-binary Go CLI tool that discovers tmux sessions running Claude Code and lets users attach with one keypress. Built with Go 1.25, Cobra (CLI framework), and Bubble Tea (TUI framework). The project has 18 Go source files across 5 packages (`cmd/`, `internal/tmux/`, `internal/process/`, `internal/profile/`, `internal/tui/`) with 3 test files containing solid unit tests.

Current eval score: **0.43** (threshold: 0.80). The project is functional but scores low because the factory's eval system has significant blind spots for Go projects.

## Critical Finding: Factory Eval Is Python-Centric

The `last_eval.json` results do NOT come from the project's `eval/score.py`. They come from the factory's built-in eval system at `remote-factory/factory/eval/` which has two modules:

- **`hygiene.py`** — tests, lint, type_check, coverage detection
- **`growth.py`** — capability_surface, observability, experiment_diversity

### What the factory eval actually does for Go:

| Dimension | Factory Behavior for Go | Current Score | Root Cause |
|---|---|---|---|
| **tests** | Runs `go test ./...` but parses output for `^ok\s+` lines. Returns 0.5 ("not detected") if no output matches. | 0.5 | Go tests exist and pass, but the factory may not be recognizing the output format correctly, OR `go test` fails silently in the eval environment (missing deps, wrong GOPATH). |
| **lint** | Only checks Python (ruff), Node (eslint), Rust (clippy). **No Go linter support.** | 0.5 | Factory does not run `go vet` or `golangci-lint`. Always returns neutral. |
| **type_check** | Only checks Python (mypy), Node (tsc). **No Go support.** | 0.5 | Go is statically typed — the compiler IS the type checker. Factory ignores this. |
| **coverage** | Only runs `pytest --cov`. **No Go support.** | 0.5 | `go test -cover` exists but factory doesn't use it. |
| **capability_surface** | **Only counts Python files.** Searches for `*.py`, uses Python AST to count functions. | 0.05 | Found 1 Python file (eval/score.py) with 4 public functions. All Go code (18 files, ~30 public functions) is invisible. |
| **observability** | Scans for `slog`, `zerolog`, `zap`, `logrus` imports AND `fmt.Printf/log.*` calls in Go files. | 0.0 | The project has zero logging statements in non-test Go files. All output uses bare `fmt.Println/Printf`. |

### Dimensions the project CAN influence:

1. **observability (weight 0.1)** — Add `log/slog` structured logging. The factory's observability eval in `growth.py` correctly scans Go files for `slog.\w+\(` patterns.
2. **tests (weight 0.15)** — Tests exist and should pass. The 0.5 score suggests they're not being detected. Improving `eval/score.py` won't help since the factory uses its own eval.
3. **capability_surface (weight 0.14)** — Cannot be fixed without modifying the factory's Python-only surface counter.

### Dimensions the project CANNOT influence (factory limitations):

- **lint** — Factory has no Go linter integration
- **type_check** — Factory has no Go type checker integration
- **coverage** — Factory has no Go coverage integration
- **capability_surface** — Factory only counts Python files

## External Research Findings

### Go Structured Logging with `log/slog`

`log/slog` (stdlib since Go 1.21) is the standard choice for structured logging in Go CLI tools. Key patterns for Cobra CLIs:

1. **Initialize in `PersistentPreRunE`**: Set up logger once in root command's persistent pre-run hook so all subcommands inherit it.
2. **Write to stderr**: Keep stdout clean for program output (e.g., `rta status --json`). Logs go to stderr.
3. **Use `--verbose` / `-v` flag**: Toggle between `slog.LevelInfo` and `slog.LevelDebug`.
4. **Use `TextHandler` for CLIs**: Human-readable format for terminal output.
5. **Key naming convention**: Use `snake_case` consistently (e.g., `session_name`, `pane_pid`).

Recommended approach for rta:
```go
// In cmd/root.go
var verbose bool

func init() {
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
}

// In PersistentPreRunE:
level := slog.LevelInfo
if verbose {
    level = slog.LevelDebug
}
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
```

Then add `slog.Debug(...)` / `slog.Info(...)` calls in key functions:
- `tmux.ListSessions()` — log session count
- `process.BuildTree()` — log process count
- `process.HasDescendant()` — log when Claude is found
- `cmd.showStatus()` — log output format
- `cmd.attachToSession()` — log match results
- `profile.InstallTo()` / `profile.RemoveFrom()` — log file operations
- `tui.refreshSessions()` — log refresh results

Sources:
- https://go.dev/blog/slog
- https://betterstack.com/community/guides/logging/logging-in-go/
- https://www.bytesizego.com/blog/cobra-cli-golang

### How the Observability Eval Scores

The factory's observability eval (in both `eval/score.py` and `growth.py`) scans Go files for logging patterns:

**Score formula** (from `eval/score.py` line 149):
```
score = 0.40 * coverage + 0.25 * has_struct + 0.20 * has_trace + 0.15 * density
```

Where:
- `coverage` = fraction of functions containing at least one log call
- `has_struct` = 1.0 if `slog`, `zerolog`, `zap`, or `logrus` is imported anywhere
- `has_trace` = 1.0 if `opentelemetry`, `trace`, or `span` appears in code
- `density` = min(1.0, total_log_calls / total_functions)

**Patterns detected as log calls** (regex):
- `\blog\.\w+\(` — standard `log` package
- `\bfmt\.Printf?\(` — fmt.Print/Printf
- `\bfmt\.Fprintf?\(` — fmt.Fprint/Fprintf
- `\bslog\.\w+\(` — structured logging

**Patterns detected as structured logging**:
- `\bslog\b`, `\bzerolog\b`, `\bzap\b`, `\blogrus\b`

## Prior Knowledge (Archive)

No archive available. The `.factory/archive/` directory is empty (all files shown as deleted in git status). Prior experiment history and source notes from previous cycles are not accessible in this worktree.

## Recommended Focus Areas

### 1. Add `log/slog` Structured Logging (HIGH IMPACT — observability: 0.0 -> ~0.6)

**Why**: The factory's observability eval correctly scans Go files for `slog.\w+\(` patterns, structured logging imports, and function-level log coverage. This is the single highest-impact improvement available.

**What**:
- Add `log/slog` import and `slog.Debug()`/`slog.Info()` calls to functions across all packages
- Add `--verbose` persistent flag to root command
- Initialize logger in `PersistentPreRunE`
- Target: log statements in 50%+ of non-trivial functions

**Projected score**: Adding slog to ~60% of functions gives: `0.40*0.6 + 0.25*1.0 + 0.20*0 + 0.15*0.5 = 0.565`

**Files to modify**: `cmd/root.go` (logger init + verbose flag), `internal/tmux/tmux.go`, `internal/process/process.go`, `internal/profile/profile.go`, `internal/profile/ssh.go`, `internal/tui/commands.go`

### 2. Ensure Test Detection Works (MEDIUM IMPACT — tests: 0.5 -> 1.0)

**Why**: The project has 3 test files with comprehensive tests that should pass. The factory reports "not detected" which means either (a) `go test ./...` fails in the eval environment, or (b) the output parser doesn't match.

**What**: The factory's `hygiene.py` runs `go test ./...` and looks for `^ok\s+` lines in stdout. Verify:
- Tests pass locally: `go test ./...` should print `ok rta/internal/tmux`, `ok rta/internal/process`, `ok rta/internal/profile`
- No network dependencies that could fail in eval environment
- No compilation errors blocking test execution

**Potential issue**: The `go.mod` specifies `go 1.25.0` — if the eval environment has an older Go version, tests would fail with a cryptic error. Consider relaxing to `go 1.22` per CLAUDE.md.

### 3. Capability Surface Is Unfixable Without Factory Changes (BLOCKED)

The factory's `growth.py` `eval_capability_surface()` only counts Python files via AST parsing. The Go project's 18 source files and ~30 public functions are invisible. This cannot be fixed by modifying the rta project.

**Current surface**: 5 (1 Python module + 4 Python functions from eval/score.py)
**Actual surface**: ~35 (4 Go packages + ~30 exported Go functions + 3 CLI entry points)

### Priority Ranking

| # | Action | Current | Target | Weight | Score Delta |
|---|--------|---------|--------|--------|-------------|
| 1 | Add slog logging | 0.0 | ~0.56 | 0.10 | +0.056 |
| 2 | Fix test detection | 0.5 | 1.0 | 0.15 | +0.075 |
| 3 | Capability surface | 0.05 | N/A | 0.14 | blocked |

**Maximum achievable score** without factory changes: ~0.56 (up from 0.43). The 0.80 threshold is not reachable because capability_surface (0.14 weight), lint (0.075), type_check (0.05), and coverage (0.125) are all capped at 0.5 or below due to missing Go support in the factory eval.

### Realistic Strategy

Focus on the two actionable items (slog + test detection) to maximize score within the factory's constraints. The combined impact is roughly +0.13 points (0.43 -> 0.56).
