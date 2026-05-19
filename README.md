# rta — Remote tmux Access

A single-binary Go CLI that discovers tmux sessions running Claude Code and lets you attach from your phone over SSH.

## Getting Started

### Prerequisites

- **Go 1.24+** (for building from source)
- **tmux** installed on the host machine
- An SSH client on your phone (Terminus, Blink Shell, Prompt 3)

### Install

```bash
# Clone and build
git clone <repo-url>
cd remote-tmux-access
make build

# Copy to your PATH
cp rta /usr/local/bin/

# Or install directly via Go
go install .
```

### Quick Setup

**1. Enable SSH on your Mac**

```bash
rta setup ssh
```

This checks whether Remote Login is enabled and walks you through fixing anything that's off. It's read-only — it won't change your system config, just tells you what to do.

**2. Auto-launch on SSH login**

```bash
rta setup
```

This adds a small block to your `.zshrc` that runs `rta` automatically when you SSH in. Local terminal sessions are unaffected. To undo:

```bash
rta setup --undo
```

**3. Connect from your phone**

Open Terminus (or any SSH app) on your iPhone and SSH into your Mac:

```
ssh youruser@your-mac-ip
```

`rta` launches automatically and shows your tmux sessions:

```
  Claude Code Sessions
  > webapp        ~/projects/webapp
    api-server    ~/projects/api

  Other tmux Sessions
    scratch       ~/tmp

  [Enter] attach  [r] refresh  [q] quit
```

Select a session and press Enter — you're in.

## Usage

### Interactive mode (default)

```bash
rta
```

Launches the TUI session picker. If only one Claude Code session exists, it attaches directly — no menu needed.

### Attach by name

```bash
rta attach webapp        # exact or substring match
rta attach web           # matches "webapp"
```

### List sessions

```bash
rta status               # table output
rta status --json        # JSON output
```

Example output:

```
SESSION      CLAUDE  DIR                  ATTACHED
webapp       *       ~/projects/webapp    yes
api-server   *       ~/projects/api
scratch              ~/tmp
```

### All commands

```
rta                          Launch TUI (or auto-attach if one Claude session)
rta attach <name>            Attach by name (substring match)
rta status                   List sessions (one per line)
rta status --json            List sessions as JSON
rta setup                    Add auto-launch to shell profile
rta setup --undo             Remove auto-launch
rta setup ssh                Check SSH server configuration
```

## How it works

1. Queries `tmux list-sessions` to find all sessions
2. For each session, gets the pane PIDs via `tmux list-panes`
3. Walks the process tree (single `ps` call) looking for `claude` descendants
4. Presents sessions in a Bubble Tea TUI, sorted by Claude detection
5. On selection, hands terminal control to `tmux attach-session` via `syscall.Exec`

Detaching (`Ctrl-b d`) returns you to the `rta` menu. The tmux session keeps running.

## Remote access

`rta` assumes direct network access (LAN, VPN, or Tailscale). For access over the public internet, layer one of:

- **Tailscale** — zero-config mesh VPN, works great with iOS
- **Cloudflare Tunnel** — expose SSH without opening ports
- **SSH port forwarding** — `ssh -R` through a jump host
