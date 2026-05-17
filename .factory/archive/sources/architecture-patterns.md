---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Architecture Patterns — Remote Terminal Access

## Pattern 1: tmux Control Mode Bridge
- tmux's `-C` flag provides a text-based protocol for programmatic access
- All pane output sent as `%output <pane_id> <data>` notifications
- Used by iTerm2 for tmux integration
- **rta decision**: Deferred to post-v1. Adds complexity without clear benefit for the attach-only use case.

## Pattern 2: PTY Wrapper (Selected by research, superseded)
- Spawn a PTY running `tmux attach`, pipe I/O over WebSocket
- Dead simple, works with any terminal application
- **rta decision**: Research recommended this for MVP, but the project pivoted away from web/WebSocket entirely. rta uses direct `tmux attach-session` via exec.

## Pattern 3: Relay Architecture (Not applicable)
- Host agent → outbound SSH/WS → relay server → client browser
- Solves NAT traversal but adds latency and operational complexity
- **rta decision**: Explicitly a non-goal. Users layer Tailscale or Cloudflare Tunnel themselves.

## Pattern Actually Used by rta: Direct Exec
- rta discovers tmux sessions, then `exec`s into `tmux attach-session -t <name>`
- No intermediate PTY, no WebSocket, no relay
- The SSH connection IS the transport layer
- Simplest possible architecture — rta is a session picker, not a proxy
