# Factory Configuration

## Goal

Build a single-binary Go CLI (rta) that discovers tmux sessions running Claude Code and lets users attach to them from iPhone SSH clients.

## Scope

### Modifiable

- main.go
- cmd/**/*.go
- internal/**/*.go
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

## Eval

### Command

```bash
python eval/score.py
```

### Threshold

0.8

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
