---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# syscall.Exec for Process Replacement — Direct Attach

## Two Attach Modes

- **`syscall.Exec`**: For `rta attach <name>` — replaces the Go process entirely with tmux. User's shell prompt reappears after tmux detach. Cleanest UX.
- **`tea.ExecProcess`**: For the TUI case — suspends TUI, runs tmux as child, resumes TUI on detach.

## syscall.Exec Implementation

```go
binary, err := exec.LookPath("tmux")  // returns full path e.g. /usr/bin/tmux
args := []string{"tmux", "attach-session", "-t", sessionName}
syscall.Exec(binary, args, os.Environ())
```

**Critical details:**
- `exec.LookPath` required — `syscall.Exec` needs an absolute path
- `args[0]` must be the program name (`"tmux"`) — Unix convention
- `os.Environ()` passes current environment
- `syscall.Exec` **only returns on error** — on success, the Go process ceases to exist
- Works on macOS and Linux; not available on Windows (not needed for rta)

## Why syscall.Exec over os/exec.Command

| | syscall.Exec | os/exec.Command |
|---|---|---|
| Process | Replaces current | Spawns child |
| After exit | Shell prompt returns | Go process resumes |
| PID | Same PID as rta | New PID |
| Use case | Final attach (no return) | Need to continue after |
