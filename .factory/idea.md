# rta — Remote tmux Access

## Vision
A single-binary Go CLI that lets you access running Claude Code sessions from an iPhone over SSH. You SSH into your Mac, `rta` launches automatically, shows you which tmux sessions have Claude Code running, and attaches you with one keypress. No web server, no wrappers, no runtime dependencies.

## Core Features

- **Claude Code Session Discovery**: Inspects the process tree of every running tmux session to find ones where the `claude` binary is active. Uses `tmux list-sessions` for session enumeration, then `tmux list-panes -t <session> -F '#{pane_pid}'` to get the shell PID per pane, and walks `/proc` (Linux) or `ps -o pid,ppid,comm` (macOS) to find `claude` in the descendant tree. Displays session name, working directory (from the pane's current path via `#{pane_current_path}`), and session dimensions.

- **Interactive TUI Menu**: A terminal UI built with Bubble Tea that lists tmux sessions containing Claude Code. Sessions with Claude Code are shown first; other tmux sessions are listed below in a secondary section. Arrow keys or `j`/`k` to navigate, Enter to attach, `r` to refresh, `q` to quit. Designed for constrained screens — renders correctly on Terminus, Blink Shell, and Prompt 3 on iOS.

- **One-Command Attach**: `rta attach <session>` attaches directly to a named session, bypassing the TUI. Supports substring matching (e.g., `rta attach web` matches a session named `webapp`). If exactly one Claude Code session exists, bare `rta` attaches to it directly — no menu needed.

- **SSH Auto-Launch**: `rta setup` appends a guarded block to `~/.zshrc` (or `~/.bashrc`) that runs `rta` only during SSH sessions. The block checks `$SSH_CONNECTION` — if set, it's an SSH session and `rta` launches; if unset, the shell starts normally. This means local terminal windows are completely unaffected. `rta setup --undo` removes the block. The injected block looks like:
  ```bash
  # rta: remote tmux access (auto-launch)
  if [ -n "$SSH_CONNECTION" ] && command -v rta >/dev/null 2>&1; then
    rta
  fi
  # end rta
  ```

- **SSH Server Configuration Helper**: `rta setup ssh` checks that macOS Remote Login is enabled, key-based auth is configured, and `~/.ssh/authorized_keys` exists. Prints clear instructions for anything that needs manual action (enabling Remote Login in System Settings, adding keys). Does not modify `sshd_config` — only inspects and advises.

- **Session Status**: `rta status` prints a non-interactive one-line-per-session summary for scripting. Shows session name, whether Claude Code is detected, working directory, and attached client count. `rta status --json` outputs structured JSON.

## Architecture

- **Language/Runtime**: Go 1.22+ — compiles to a single static binary with zero runtime dependencies. The research identifies Go as the production-grade choice (Option B), proven by Upterm. Distribution is a single file copy or `go install`.

- **Framework**: No web framework. This is a pure CLI/TUI application. Uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Charm ecosystem) for the interactive session picker — it's the standard Go TUI framework, handles terminal resize and input correctly, and renders well on constrained mobile terminals.

- **Data Storage**: None. tmux is the source of truth. No database, no config files, no state on disk. Session metadata is queried live on every invocation.

- **Key Libraries**:
  - `github.com/charmbracelet/bubbletea` — TUI framework (Elm-architecture, handles resize, input, rendering)
  - `github.com/charmbracelet/lipgloss` — TUI styling (borders, colors, layout)
  - `github.com/spf13/cobra` — CLI subcommand routing and flag parsing
  - `os/exec` — all tmux interaction via `tmux` CLI commands (list-sessions, list-panes, attach-session)
  - `os/exec` — process tree inspection via `ps` commands on macOS, `/proc` traversal on Linux

- **tmux Interaction**: All tmux operations go through the `tmux` CLI binary — no libraries, no control mode, no sockets. `rta` calls `tmux list-sessions`, `tmux list-panes`, and `tmux attach-session` as subprocesses. This is the simplest and most robust approach — tmux CLI is stable and well-documented.

- **Process Detection**: To determine if a tmux pane is running Claude Code, `rta` gets the pane's shell PID via tmux format strings, then walks the process tree looking for a process named `claude`. On macOS this uses `ps -o pid,ppid,comm`; on Linux it reads `/proc/<pid>/stat`. The walk is shallow (typically 3-4 levels: tmux → shell → claude) and fast.

## User Interface

**Primary flow (from iPhone with Terminus):**

1. Claude Code sessions are already running in tmux on the Mac (started locally earlier)
2. User opens Terminus on iPhone, connects via SSH (local network or Tailscale)
3. `rta` launches automatically (configured via `rta setup`)
4. TUI shows active sessions:
   ```
   ┌─ rta ─────────────────────────────────────┐
   │                                            │
   │  Claude Code Sessions                      │
   │  ▸ webapp        ~/projects/webapp          │
   │    api-server    ~/projects/api              │
   │                                            │
   │  Other tmux Sessions                       │
   │    scratch       ~/tmp                      │
   │                                            │
   │  [Enter] attach  [r] refresh  [q] quit     │
   └────────────────────────────────────────────┘
   ```
5. User selects a session and presses Enter
6. `rta` execs into `tmux attach-session -t <name>` — the TUI is gone, they're in the tmux session
7. Detaching (`Ctrl-b d`) returns to the `rta` menu

**CLI interface:**

```
rta                          # Launch TUI (or auto-attach if exactly one Claude session)
rta attach <name>            # Attach by name (substring match)
rta status                   # List sessions (one line per session)
rta status --json            # List sessions as JSON
rta setup                    # Add auto-launch block to shell profile
rta setup --undo             # Remove auto-launch block
rta setup ssh                # Check/guide SSH server configuration
```

## Non-Goals (v1)

- **Web UI** — no browser-based access; native SSH terminal apps only
- **Multi-agent support** — v1 detects Claude Code only. Aider, Cursor, Copilot detection is a future addition
- **Wrapping or modifying Claude Code** — `rta` never starts, stops, or sends input to Claude Code. It only discovers and attaches to tmux sessions where Claude Code is already running
- **Custom SSH server** — uses the system's OpenSSH. No custom SSH implementation
- **NAT traversal / tunneling** — assumes direct network access (LAN or VPN). Users layer Tailscale or Cloudflare Tunnel themselves
- **Mosh integration** — better mobile connectivity but adds install complexity. Deferred
- **Multi-user access control** — single-user tool (your laptop, your account)
- **Session recording / playback**
- **Windows support** — macOS and Linux only
- **tmux control mode** — v1 uses simple `tmux attach`. Control mode (`-C`) adds complexity without clear benefit for the attach-only use case

## Open Questions

- **Installation method**: Should `rta` be distributed via Homebrew tap, GitHub releases with prebuilt binaries, or just `go install`? All three are feasible for a single Go binary. Homebrew tap is the most iPhone-user-friendly since the target audience is macOS.

## Changes from Prior Draft

1. **Go instead of Python** — replaced the entire Python stack (textual, libtmux, psutil, click) with Go equivalents (Bubble Tea, Lip Gloss, Cobra, os/exec). Single binary, zero runtime dependencies.
2. **Claude Code only** — removed generic agent detection (Aider, Cursor). v1 focuses exclusively on finding the `claude` binary in process trees. Other agents moved to Non-Goals.
3. **Explicit auto-launch mechanism** — documented exactly what `rta setup` injects into the shell profile: a guarded block checking `$SSH_CONNECTION` before launching `rta`. Showed the literal shell code.
4. **Clarified non-wrapper stance** — made it explicit throughout the spec that `rta` never starts, stops, wraps, or modifies Claude Code. It's purely a discovery and attachment layer on top of tmux.
5. **Removed session naming convention dependency** — detection is purely process-tree based. No assumptions about how tmux sessions are named.
