# rta — Remote tmux Access

Single-binary Go CLI that discovers tmux sessions running Claude Code and attaches with one keypress.

## Build

- Go 1.22+
- `make build` to build
- `make test` to test
- `make lint` to lint

## Conventions

- Run `go vet ./...` before committing
- Run `go test ./...` before opening a PR
- Keep packages small and focused: `internal/tmux/`, `internal/process/`, `internal/tui/`, `internal/profile/`
- Use `cmd/` for cobra command definitions
