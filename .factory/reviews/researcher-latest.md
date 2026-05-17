# Researcher Agent Output

- **timestamp:** 2026-05-17T18:06:03Z
- **exit_code:** 0

---

Research complete. Wrote `.factory/strategy/research.md` covering all six topics with concrete Go code patterns:

1. **Bubble Tea** — custom list selector with cursor/j/k, `tea.ExecProcess` for suspending TUI and handing terminal to tmux, lipgloss styling
2. **Process tree** — single `ps -ax -o pid,ppid,comm` call, parse into parent→children map, BFS walk from pane PID looking for `claude`
3. **tmux parsing** — `list-sessions`/`list-panes` with `-F` format strings using pipe delimiter, error handling for no-server/no-sessions cases, note that `session_width/height` were removed in tmux 2.9
4. **Cobra** — root command runs TUI via `RunE`, positional args for attach, `BoolVarP` for `--json`/`--undo`, nested `setup ssh` subcommand
5. **syscall.Exec** — replaces process entirely for `rta attach`, `tea.ExecProcess` for TUI case; includes `exec.LookPath` for finding tmux binary
6. **Profile injection** — marker-based block with idempotent install/remove, shell detection via `$SHELL`, permission preservation
