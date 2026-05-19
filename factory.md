# Factory Configuration

## Goal

Build a single-binary Go CLI (rta) that discovers tmux sessions running Claude Code, presents them in a TUI with recap descriptions, and lets users attach with one keypress. Includes a Claude Code skill (`/enable-rta`) to onboard non-tmux sessions.

## Scope

### Modifiable

- main.go
- cmd/**/*.go
- internal/**/*.go
- skill/**/*
- eval/score.py
- Makefile
- CLAUDE.md

### Read-only

- go.mod
- go.sum
- factory.md

## Guards

- Do not delete or overwrite existing tests
- Do not modify files outside the declared scope
- Do not introduce secrets or credentials into the repository
- All code must have test coverage

## Eval

### Command

```bash
python3 eval/score.py
```

### Threshold

0.80

## Target Branch

main

## Smoke Test

```bash
go build -o ./rta . && ./rta status 2>&1 | grep -q "SESSION\|No tmux"
```

## Constraints

- Prefer small, incremental changes over large rewrites
- Each change should be accompanied by at least one test
- Follow the existing code style and conventions
- All tmux interaction via CLI commands (no sockets, no control mode)
- Single binary, zero runtime dependencies
