---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Process Tree Walking on macOS — Claude Detection

## ps Command

```bash
ps -ax -o pid,ppid,comm
```

- `ps -ax` and `ps -e` are equivalent on macOS (all processes)
- `comm` returns the base executable name (e.g., `claude`, `bash`), not full path — ideal for matching
- Header line must be skipped when parsing
- Output is space-separated with leading whitespace

## Design: Single ps Call + BFS Walk

Build the full process tree once per refresh (single `ps` call), then walk it for each pane. Avoids N+1 subprocess calls. Data structures:
- `children map[int][]int` — parent -> children mapping
- `names map[int]string` — pid -> command name mapping

Given a pane's shell PID from tmux, BFS all descendants looking for process named `claude`.

## Edge Cases

- Zombie processes (`Z` state) still appear in ps — won't cause false positives since they'd only match if named `claude`
- Process may exit between `ps` call and use — not a problem since we only read from the snapshot
- macOS background processes reparent to launchd (PID 1) — irrelevant since we walk DOWN from known pane PIDs, not UP
- Linux alternative: read `/proc/<pid>/stat` instead of `ps` — both paths must be implemented
