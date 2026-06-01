---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 2
verdict: REVERT
score_delta: "+0.075 (0.426 → 0.5011, still below 0.80)"
date: 2026-05-19
source: factory-archivist
---

# Experiment #2 (Optimization Cycle): REVERT — Add slog structured logging

## Hypothesis
Add log/slog structured logging across all packages with a --verbose flag to improve observability score.

## Result
**REVERT** — Score improved from 0.426 to 0.5011 (+0.075), but still far below threshold 0.80. Precheck failed on `score_direction` (improvement insufficient) and `scope`.

This is a **systemic issue**, not a code quality issue. The factory eval is Python-centric and structurally cannot score Go projects above ~0.56.

## Score Breakdown (post-change)
| Dimension | Score | Weight | Notes |
|-----------|-------|--------|-------|
| tests | 1.0 | 0.15 | Go tests pass |
| lint | 0.5 | 0.075 | No Go linter detected |
| type_check | 0.5 | 0.05 | No Go type checker detected |
| coverage | 0.5 | 0.125 | No Go coverage tool detected |
| guard_patterns | 0.75 | 0.05 | 6/8 passed |
| config_parser | 1.0 | 0.05 | OK |
| capability_surface | 0.05 | 0.14 | Python AST only — 5/100 surface |
| experiment_diversity | 0.5 | 0.11 | Too few experiments |
| observability | 0.0 | 0.10 | Still 0 despite slog addition |
| research_grounding | 0.52 | 0.08 | Sources present |
| factory_effectiveness | 0.5 | 0.07 | Too few experiments |

## Why Reverted
The code change itself was well-executed (CEO review CLEAN on all 7 checks, PR #4 merged). However:
1. Score 0.5011 < threshold 0.80 triggers automatic revert
2. The threshold is unreachable for Go projects — blocked dimensions (capability_surface, lint, type_check, coverage) have combined weight ~0.40
3. Observability scored 0.0 even after adding slog — the eval likely checks Python-specific observability patterns
4. Maximum achievable score is ~0.56 regardless of code quality

## What Was Built (then reverted)
- 8 files modified, +47 lines, zero new dependencies
- `log/slog` with TextHandler to stderr, `--verbose` flag
- See [experiment note 010](remote-tmux-access-010.md) for full change details

## Systemic Analysis
The factory eval system cannot meaningfully score Go projects:
- `capability_surface`: Counts Python AST nodes only
- `lint`: Detects ruff/eslint/clippy, not golangci-lint or go vet
- `type_check`: Detects mypy/pyright, not Go's compile-time type checking
- `coverage`: Detects pytest coverage, not go test -cover
- `observability`: Likely checks Python logging patterns, not Go slog

This means ALL future experiments on this Go project will be reverted unless the eval system gains Go support or the threshold is lowered.

## Links
- Project: remote-tmux-access
- Issue: #3
- PR: #4 (merged then reverted)
- Prior note (KEEP verdict before revert): [experiment 010](remote-tmux-access-010.md)
- Related pattern: [Factory Eval Python Bias](../patterns/patterns.md#factory-eval-python-bias-creates-hard-score-ceiling-for-non-python-projects)
