---
tags:
  - factory
  - strategy
  - remote-tmux-access
date: 2026-05-17
source: factory-archivist
---

# Strategy: remote-tmux-access — 2026-05-17

## Ideation Process

### Phase 1: Research (04:19–04:24 UTC)

The researcher agent surveyed the landscape of remote terminal sharing tools. Six existing projects were analyzed:

1. **tmate** — SSH-based terminal sharing, fork of tmux 2.x. Battle-tested but permanently behind upstream tmux. Written in C, hard to extend.
2. **sshx** — Rust-based with web UI, E2E encryption, collaborative features. Complex self-hosting (gRPC, Redis, mesh networking).
3. **Upterm** — Go-based, reverse SSH tunnels, works with any command. Closest architectural precedent for this project.
4. **ttyd** — C tool exposing terminal via WebSocket + xterm.js. Simple but no NAT traversal or encryption.
5. **GoTTY** — Go CLI-to-web tool, largely unmaintained.
6. **pyxtermjs** — Minimal Python reference implementation, not production-grade.

**Key finding**: None of these focus on accessing *existing* tmux sessions. All either fork tmux or spawn their own PTY. The unique value proposition is bridging to real, running tmux sessions.

### Phase 2: Initial Spec (04:25 UTC)

The distiller produced a Python-based spec following the researcher's recommendation:
- Python + asyncio + websockets + xterm.js
- PTY bridge pattern (spawn `tmux attach` in a PTY, pipe I/O over WebSocket)
- Token-based auth, resize handling, clean disconnect

CEO verdict: PROCEED — spec captures intent, technology justified by research, scope tight.

### Phase 3: User Feedback & Pivot (04:35 UTC)

User redirected the project significantly:
- **Target platform**: iOS via SSH terminal apps (Terminus, Blink Shell, Prompt 3), not web browsers
- **No web UI**: Pure CLI/TUI, SSH-native access
- **Focus**: Claude Code session discovery specifically, not generic terminal sharing

The distiller refined the spec to reflect iOS/SSH focus while keeping the Python stack.

### Phase 4: Go Rewrite (13:00 UTC)

User directed a complete technology pivot:
1. **Go instead of Python** — single binary, zero runtime dependencies
2. **Claude Code only** — detect `claude` binary in process trees, not generic agents
3. **Explicit auto-launch** — documented exact shell block for `rta setup` with `$SSH_CONNECTION` guard
4. **Non-wrapper stance** — rta never starts/stops/modifies Claude Code

The final spec was approved as `.factory/strategy/current.md`.

## Strategic Rationale

The project fills a genuine personal workflow gap: the user runs multiple Claude Code sessions in tmux on their Mac and wants to check on them from their iPhone. Existing tools are over-engineered for this use case (relay servers, web UIs, PTY spawning). The simplest solution is a Go binary that discovers tmux sessions and attaches to them over plain SSH.

## Risk Assessment

- **Low technical risk** — tmux CLI is stable, process tree walking is well-understood, Bubble Tea is mature
- **Low scope risk** — narrow feature set (discover, list, attach) with clear non-goals
- **Platform risk** — iOS SSH app rendering quality varies; TUI must be tested on constrained screens
