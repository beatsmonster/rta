---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Research Pitfalls — Applicability to rta

The research identified 8 pitfalls for remote terminal access. Since rta pivoted from web-based to SSH-native, most are no longer applicable, but some inform related concerns.

## Still Relevant

### Terminal Resize on Constrained Screens
- **Original concern**: Browser resize events trigger PTY resize storms
- **rta version**: iOS SSH apps have fixed screen sizes, but rotating the phone or split-screen on iPad could trigger resize. tmux handles resize natively, but TUI (Bubble Tea) must render correctly at various dimensions.
- **Mitigation**: Test on Terminus, Blink Shell, and Prompt 3 at various orientations.

### Platform Differences (macOS vs Linux)
- **Original concern**: PTY handling differs between platforms
- **rta version**: Process tree walking differs — `ps -o pid,ppid,comm` on macOS vs `/proc/<pid>/stat` on Linux. Both paths must be implemented and tested.

### tmux Session Discovery
- **Original concern**: Server must run as same user who owns tmux sessions
- **rta version**: Directly applicable. rta must be run as the user who owns the tmux server. SSH handles this naturally (you SSH in as yourself).

## No Longer Applicable (due to SSH-native pivot)

- **WebSocket reconnection** — SSH handles connection persistence
- **Encoding/binary data** — SSH terminal handles encoding
- **Security (WSS, tokens)** — SSH key-based auth is the security layer
- **Latency/predictive echo** — SSH is already low-latency for terminal use
- **PTY blocking** — rta execs into tmux, doesn't proxy PTY I/O
