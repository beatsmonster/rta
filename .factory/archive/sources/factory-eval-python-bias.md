---
tags:
  - factory
  - source
source: factory-archivist
date: 2026-05-19
---

# Factory Eval System Is Python-Centric

## Finding

The factory's built-in eval system (`hygiene.py` + `growth.py`) was designed for Python projects. Key limitations for Go projects:

### capability_surface (weight 0.14) — BLOCKED
- Only counts Python files (`*.py`) using Python AST parsing
- All 18 Go source files and ~30 exported functions are invisible
- Current score: 0.05 (found 1 Python file `eval/score.py` with 4 functions)
- Cannot be fixed without modifying the factory eval code

### lint — No Go Support
- Only checks ruff (Python), eslint (Node), clippy (Rust)
- `go vet` and `golangci-lint` are not recognized

### type_check — No Go Support
- Only checks mypy (Python), tsc (Node)
- Go is statically typed at compile time; factory ignores this

### coverage — No Go Support
- Only runs `pytest --cov`
- `go test -cover` exists but is not used

## Impact on Score Ceiling

Maximum achievable score without factory changes: ~0.56 (threshold: 0.80). The gap is structural, not quality-related.
