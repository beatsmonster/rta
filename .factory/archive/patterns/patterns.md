---
tags:
  - factory
  - patterns
source: factory-archivist
---

# Cross-Project Patterns

## Research-to-Spec Pivot
Discovered in remote-tmux-access ideation (2026-05-17).
Research recommended a Python + WebSocket + xterm.js web-based architecture. User feedback pivoted the project to Go + SSH-native CLI/TUI with zero web components. The research was still valuable — it identified the key differentiator (bridging to existing tmux sessions) and mapped the competitive landscape — but the technology recommendations were entirely replaced. **Lesson**: Research should emphasize *what* to build and *why* before *how* to build it. Technology recommendations are the most likely part to be overridden by user preferences.

## Implementation Research as Executable Documentation
Discovered in remote-tmux-access implementation research (2026-05-17).
When research produces code samples for a specific tech stack (Go + Bubble Tea + Cobra + tmux), the samples should be directly copy-pasteable — complete imports, correct types, real API calls. The researcher's implementation findings for rta were immediately usable because they included exact function signatures, message types, and edge case handling. **Lesson**: Post-pivot implementation research should produce code-level patterns, not just architectural descriptions. The strategist can then sequence them without re-researching APIs.

## Platform-Specific Error Strings Require Real-Device Testing
Discovered in remote-tmux-access Phase 4 (experiment #4, 2026-05-17).
tmux on macOS reports "error connecting to /tmp/tmux-xxx/default" when no server is running, while Linux reports "no server running on /tmp/tmux-xxx/default". The research source notes warned about platform differences, but the exact error string divergence was only caught during real-device testing. **Lesson**: When wrapping CLI tools that produce human-readable error messages, test error paths on all target platforms. Error string matching is fragile — check for multiple known variants or use exit codes where possible.

## Incremental 8-Phase Build With Zero Reverts
Discovered in remote-tmux-access full build (2026-05-17).
An 8-phase incremental build plan achieved a 100% keep rate (8/8 phases kept, 0 reverts). Key factors: (1) each phase had a clear dependency on prior phases, preventing integration surprises; (2) deferred optional dependencies (Charm libs) until the phase that needed them; (3) platform-specific bugs were caught and fixed within the phase that introduced them rather than accumulating. **Lesson**: For greenfield CLI projects, a dependency-ordered phase plan with one feature per phase and immediate test/lint gates produces reliable incremental progress. The keep rate tracks build quality — 100% means the strategy's phase decomposition matched the code's actual dependency graph.
