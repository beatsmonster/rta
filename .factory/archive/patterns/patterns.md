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
