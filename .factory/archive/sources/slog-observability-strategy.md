---
tags:
  - factory
  - source
source: factory-archivist
date: 2026-05-19
---

# Observability via log/slog — Primary Actionable Dimension

## Finding

Observability (currently 0.0, weight 0.10) is the most impactful dimension the project can improve. The factory's observability eval scans Go source files directly (no Go binary needed), looking for:

### Score Formula
```
score = 0.40 * coverage + 0.25 * has_struct + 0.20 * has_trace + 0.15 * density
```

Where:
- `coverage` = fraction of functions with at least one log call
- `has_struct` = 1.0 if `slog`, `zerolog`, `zap`, or `logrus` is imported
- `has_trace` = 1.0 if `opentelemetry`, `trace`, or `span` appears
- `density` = min(1.0, total_log_calls / total_functions)

### Detected Patterns (regex)
- `\blog\.\w+\(` — standard log package
- `\bfmt\.Printf?\(` — fmt.Print/Printf
- `\bslog\.\w+\(` — structured logging

### Recommended Approach
Use `log/slog` (stdlib since Go 1.21):
- Initialize in root command's `PersistentPreRunE` with `--verbose` flag
- Write to stderr (keep stdout clean for program output)
- Add `slog.Debug()`/`slog.Info()` to 50%+ of functions
- Projected score: ~0.56 (from `0.40*0.6 + 0.25*1.0 + 0.20*0 + 0.15*0.5`)

### Key Files to Modify
`cmd/root.go`, `internal/tmux/tmux.go`, `internal/process/process.go`, `internal/profile/profile.go`, `internal/profile/ssh.go`, `internal/tui/commands.go`

### Sources
- https://go.dev/blog/slog
- https://betterstack.com/community/guides/logging/logging-in-go/
