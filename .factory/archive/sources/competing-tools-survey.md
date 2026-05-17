---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Competing Tools Survey — Remote Terminal Sharing

## tmate — SSH-based Terminal Sharing
- **URL**: https://tmate.io | https://github.com/tmate-io/tmate-ssh-server
- **Architecture**: Fork of tmux 2.x → SSH tunnel via libssh → relay server → remote SSH clients
- **Key detail**: Permanently locked to tmux 2.x codebase (C). Geo-distributed relays. Session isolation via PID namespaces.
- **Relevance to rta**: Demonstrates demand for remote tmux access but takes the wrong approach (forking tmux). rta bridges to real tmux instead.

## sshx — Collaborative Live Terminal
- **URL**: https://sshx.io | https://github.com/ekzhang/sshx
- **Architecture**: Rust CLI → gRPC → relay (Rust/Axum/Tonic) → WebSocket → browser (Svelte + xterm.js). Redis for coordination.
- **Key detail**: E2E encryption (Argon2 + AES), Mosh-style predictive echo, infinite canvas UI. Deployed on Fly.io.
- **Relevance to rta**: Shows state-of-the-art in terminal sharing UX, but is web-first and spawns its own PTY. Over-engineered for rta's use case.

## Upterm — Go-based Terminal Sharing
- **URL**: https://upterm.dev | https://github.com/owenthereal/upterm
- **Architecture**: Go. Host SSH server + reverse tunnel to `uptermd` relay. Supports SSH and SSH-over-WebSocket.
- **Key detail**: Closest architectural precedent — single Go binary, SSH-native. Proves Go is viable for this category.
- **Relevance to rta**: Validates Go as the right language choice. However, Upterm still spawns its own PTY rather than attaching to existing tmux sessions.

## ttyd — Terminal Over Web
- **URL**: https://github.com/nicm/ttyd
- **Architecture**: C. Spawns PTY → WebSocket server → browser (xterm.js).
- **Relevance to rta**: Simple reference for PTY-to-browser bridge. Not applicable since rta skips the web layer entirely.

## GoTTY — CLI to Web (Go)
- **URL**: https://github.com/yudai/gotty
- **Key detail**: Largely unmaintained. No NAT traversal, limited security.
- **Relevance to rta**: Minimal — validates that Go + terminal tools is a well-trodden path.

## pyxtermjs — Python Reference
- **URL**: https://github.com/cs01/pyxtermjs
- **Key detail**: Flask + SocketIO + pty module + xterm.js. Proof-of-concept quality.
- **Relevance to rta**: Initially informed the Python MVP recommendation, but the project pivoted to Go.
