---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# tmux CLI Output Parsing — Format Strings and Error Handling

## list-sessions

```bash
tmux list-sessions -F '#{session_name}|#{session_attached}|#{session_windows}'
```

Available format variables: `session_name`, `session_attached` (1/0), `session_windows`.

**IMPORTANT:** `#{session_width}` and `#{session_height}` were removed in tmux 2.9. Use `#{window_width}` / `#{window_height}` or `#{pane_width}` / `#{pane_height}` instead.

## list-panes

```bash
tmux list-panes -t <session> -F '#{pane_pid}|#{pane_current_path}'
```

- `#{pane_pid}` returns the shell PID (typically bash/zsh). Claude runs as a descendant.
- `#{pane_current_path}` returns the pane's current working directory.

## Delimiter Choice

**Use pipe `|` instead of colon `:`**. File paths contain colons on some systems, and `|` is never valid in paths. Avoids ambiguous parsing. Parse with `strings.SplitN(line, "|", N)`.

## Error Handling

| Scenario | Exit Code | Stderr |
|---|---|---|
| No tmux server running | 1 | `no server running on ...` |
| Server running, no sessions | 1 | `no sessions` |
| Session not found | 1 | `can't find session <name>` |
| Success | 0 | (empty) |

Both "no server" and "no sessions" return exit code 1. Parse stderr to distinguish if needed, or treat both as "no sessions available."
