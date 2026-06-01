---
tags:
  - factory
  - source
source: factory-archivist
date: 2026-05-19
---

# Go Not Installed in Eval Environment

## Finding

The eval environment does not have Go installed. Commands `go test ./...`, `go vet ./...`, and `go build` all fail with "command not found". This means the factory's built-in `hygiene.py` eval cannot detect tests, lint, or type checking for this Go project — all return 0.5 ("not detected").

## Impact

- **tests** (weight 0.15): Returns 0.5 despite 25 passing tests in 3 test files
- **lint** (weight 0.075): Returns 0.5 — factory has no Go linter support anyway
- **type_check** (weight 0.05): Returns 0.5 — Go's compiler is the type checker, factory ignores this
- **coverage** (weight 0.125): Returns 0.5 — factory only runs `pytest --cov`

Combined weight of blocked dimensions: 0.40 (out of 1.0). This represents a hard ceiling on achievable score.

## Implication

Hypotheses must focus on dimensions that score via source scanning (no Go binary needed), not on dimensions that require running Go commands.
