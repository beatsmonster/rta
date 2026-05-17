---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Cobra CLI Subcommand Patterns — rta Project Layout

## Recommended Directory Structure

```
rta/
├── main.go           # calls cmd.Execute()
├── cmd/
│   ├── root.go       # root command (launches TUI)
│   ├── attach.go     # rta attach <name>
│   ├── status.go     # rta status [--json]
│   ├── setup.go      # rta setup [--undo]
│   └── setup_ssh.go  # rta setup ssh
└── internal/
    ├── tmux/         # tmux interaction
    ├── process/      # process tree walking
    ├── tui/          # Bubble Tea TUI
    └── profile/      # shell profile injection
```

## Key Patterns

- **Root command runs TUI**: `rootCmd.RunE` launches the Bubble Tea TUI when no subcommand given
- **Positional args**: `attach` uses `cobra.ExactArgs(1)` with `args[0]` for session name (substring match)
- **Boolean flags**: `--json` on status, `--undo` on setup, registered via `Flags().BoolVarP()`
- **Nested subcommands**: `setup ssh` via `setupCmd.AddCommand(setupSSHCmd)` in init()
- **main.go**: Minimal — just `package main; import "rta/cmd"; func main() { cmd.Execute() }`
