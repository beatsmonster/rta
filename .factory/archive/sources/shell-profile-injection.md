---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Shell Profile Block Injection/Removal — Auto-Launch Setup

## The Auto-Launch Block

```bash
# rta: remote tmux access (auto-launch)
if [ -n "$SSH_CONNECTION" ] && command -v rta >/dev/null 2>&1; then
  rta
fi
# end rta
```

Only triggers when:
1. `$SSH_CONNECTION` is set (SSH session, not local terminal)
2. `rta` is on PATH

## Shell Detection

Reads `$SHELL` env var; matches suffix `zsh` -> `~/.zshrc`, `bash` -> `~/.bashrc`. Returns error for unsupported shells.

## Injection Safety

- **Idempotent**: Checks for `startMarker` before injecting; no-op if already present
- **File creation**: `os.O_CREATE` handles missing `.zshrc` / `.bashrc`
- **Permission preservation**: Read existing permissions with `os.Stat()`, apply with `os.WriteFile()`

## Removal

Line-by-line scan: set `removing = true` between start/end markers, skip those lines, write back remaining content. Preserves surrounding content and file permissions exactly.

## Marker-Based Approach

Uses `# rta: remote tmux access (auto-launch)` and `# end rta` as delimiters. This pattern is robust — grep-friendly, human-readable, and supports idempotent install/uninstall cycles.
