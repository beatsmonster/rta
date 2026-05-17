---
tags:
  - factory
  - strategy
  - remote-tmux-access
date: 2026-05-17
source: factory-archivist
---

# Strategy: remote-tmux-access — Build Plan (CEO-Approved 2026-05-17)

## CEO Verdict

**PROCEED** — 8-phase plan is well-structured, dependency-ordered, and achievable. Each phase is one PR's worth of work. Phase 1 is correctly scaffold + eval. Implementation details reference research code samples directly. Dependency graph is correct — phases 7/8 can parallelize after 6.

No issues found. No items deferred. No human intervention needed.

## Approved Build Phases

### Phase 1: Project Scaffold + Go Module + Eval Harness
- go.mod (module `rta`, Go 1.22+, bubbletea v2, lipgloss v2, cobra)
- main.go, cmd/root.go, internal/{tmux,process,tui,profile} stubs
- Makefile (build, test, lint, install), CLAUDE.md, factory.md, .gitignore
- Eval: builds, tests_pass, vet_clean, binary_runs

### Phase 2: tmux Parsing (`internal/tmux/`)
- Session/Pane structs, ListSessions, ListPanes, AttachSession
- Pipe `|` delimiter (not colon), SplitN, "no server" vs "no sessions" error handling
- Unit tests with mock output

### Phase 3: Process Tree (`internal/process/`)
- BuildTree from single `ps -ax -o pid,ppid,comm` call
- BFS HasDescendant walk, DetectClaude convenience
- Unit tests with mock ps output

### Phase 4: Cobra Commands (`cmd/`)
- root (launch TUI), attach (substring match + syscall.Exec), status (--json), setup (placeholder)
- Process tree built once, walked per pane

### Phase 5: syscall.Exec Attach (End-to-End)
- Integration testing: ambiguous match, no sessions, stale sessions
- FindSession substring matcher in internal/tmux/

### Phase 6: Bubble Tea TUI (`internal/tui/`)
- Two-section list: Claude Code Sessions / Other tmux Sessions
- j/k/arrows navigate, Enter attaches via tea.ExecProcess, r refreshes, q quits
- Auto-attach when exactly one Claude session exists

### Phase 7: Shell Profile Injection (`internal/profile/`)
- Marker-based idempotent install/remove
- `$SSH_CONNECTION` guard, shell detection from `$SHELL`
- Tests for empty file, idempotent, remove, no-op, permissions

### Phase 8: SSH Config Checker
- Read-only inspection: sshd running, authorized_keys, permissions
- Pass/fail output with remediation instructions
- Never modifies system files

## Dependency Graph

```
1 → 2 → 3 → 4 → 5 → 6 → {7, 8}
```

Phases 7 and 8 are independent after Phase 6.

## Key Technical Decisions

- **Pipe delimiter** for tmux format strings (colons appear in paths)
- **Single `ps` call** + BFS walk (avoid N+1 subprocess calls)
- **`syscall.Exec`** for attach (replaces Go process with tmux)
- **`tea.ExecProcess`** for TUI→tmux handoff (suspends TUI, resumes on detach)
- **Marker-based profile injection** (safe, idempotent, cleanly reversible)
- **Custom list over bubbles/list** (simpler for two-section display)
