---
tags:
  - factory
  - source
  - remote-tmux-access
source: factory-archivist
date: 2026-05-17
---

# Bubble Tea TUI Patterns — Implementation Reference

## Model-View-Update Architecture

All Bubble Tea apps follow The Elm Architecture: Model (state struct) -> Init (initial commands) -> Update (handle messages, return commands) -> View (render output).

## Custom List Selector (Recommended over bubbles/list)

A custom list is simpler for rta's use case — only needs cursor movement, two sections (Claude sessions / other sessions), and Enter to select. The model struct holds `sessions []Session`, `cursor int`, `width int`, `height int`.

## Key Message Types

- `tea.KeyPressMsg` — use `msg.String()` for `"up"`, `"down"`, `"enter"`, `"j"`, `"k"`, `"q"`, `"ctrl+c"`
- `tea.WindowSizeMsg` — has `Width` and `Height` fields, sent at startup and on every resize
- Custom `tmuxFinishedMsg{err error}` — returned when tmux detaches

## tea.ExecProcess — Terminal Handoff to tmux

Critical pattern for TUI -> tmux. When user selects a session:

1. Return `tea.ExecProcess(cmd, callback)` from `Update()`
2. Bubble Tea restores terminal to normal mode, releases stdin/stdout/stderr to tmux
3. When user detaches from tmux (Ctrl-b d), tmux exits, callback fires
4. Resulting message delivered to next `Update()` call — TUI re-renders

```go
func attachTmux(sessionName string) tea.Cmd {
    c := exec.Command("tmux", "attach-session", "-t", sessionName)
    return tea.ExecProcess(c, func(err error) tea.Msg {
        return tmuxFinishedMsg{err}
    })
}
```

## Lipgloss Styling

Use `lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))` and `style.Render("text")` in `View()`.
