---
tags:
  - factory
  - experiment
  - remote-tmux-access
project: remote-tmux-access
experiment_id: 5
verdict: keep
score_delta: null
date: 2026-05-19
source: factory-archivist
---

# Experiment #5: /enable-rta skill + setup skill subcommand

## Hypothesis
/enable-rta Claude Code skill and `rta setup skill` subcommand for installing the skill.

## Result
**KEEP** — Skill file created with session migration instructions, Go subcommand uses go:embed. PR #6.

## What Changed
- Created `skill/enable-rta/SKILL.md` with session migration instructions for Claude Code
- Added `rta setup skill` subcommand that installs the skill file to `~/.claude/skills/`
- Uses `go:embed` to bundle the skill file into the binary
- Tests at 91% cmd coverage
- Review pipeline: full, 1 iteration

## Notes
- The skill teaches Claude Code how to use rta for remote access
- go:embed ensures the skill file ships with the single binary — no external dependencies

## Links
- Project: remote-tmux-access
- Issue: #8
- PR: #6
