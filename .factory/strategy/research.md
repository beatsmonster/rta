# Research — rta Implementation

## Project Summary

`rta` is a single-binary Go CLI that discovers tmux sessions running Claude Code, presents them in a TUI, and attaches with one keypress. No web server, no runtime dependencies. Target: iPhone SSH clients (Terminus, Blink Shell, Prompt 3).

---

## 1. Bubble Tea TUI — List Selector & Exec

### Model-View-Update Architecture

All Bubble Tea apps follow The Elm Architecture: Model (state struct) → Init (initial commands) → Update (handle messages, return commands) → View (render output).

### Custom List Selector (Recommended over bubbles/list)

A custom list is simpler for this use case — we only need cursor movement, two sections (Claude sessions / other sessions), and Enter to select.

```go
type model struct {
    sessions []Session
    cursor   int
    width    int
    height   int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 { m.cursor-- }
        case "down", "j":
            if m.cursor < len(m.sessions)-1 { m.cursor++ }
        case "enter":
            return m, attachTmux(m.sessions[m.cursor].Name)
        case "r":
            return m, refreshSessions()
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}
```

**Key message types:**
- `tea.KeyPressMsg` — use `msg.String()` for `"up"`, `"down"`, `"enter"`, `"j"`, `"k"`, `"q"`, `"ctrl+c"`
- `tea.WindowSizeMsg` — has `Width` and `Height` fields, sent at startup and on every resize

### tea.ExecProcess — Handing Terminal to tmux

This is the critical pattern. When the user selects a session from the TUI, Bubble Tea suspends itself, gives full terminal control to tmux, and resumes when tmux detaches.

```go
type tmuxFinishedMsg struct{ err error }

func attachTmux(sessionName string) tea.Cmd {
    c := exec.Command("tmux", "attach-session", "-t", sessionName)
    return tea.ExecProcess(c, func(err error) tea.Msg {
        return tmuxFinishedMsg{err}
    })
}
```

**How it works:**
1. Returning `tea.ExecProcess(cmd, callback)` from `Update()` causes Bubble Tea to restore the terminal to normal mode, release stdin/stdout/stderr to the spawned process, and block until it exits.
2. When the user detaches from tmux (Ctrl-b d), the tmux process exits, the callback fires, and the resulting message is delivered to the next `Update()` call.
3. The TUI re-renders — the user is back at the session list.

### Lipgloss Styling

```go
title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).MarginBottom(1)
selected := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Background(lipgloss.Color("25")).PaddingLeft(2)
normal := lipgloss.NewStyle().PaddingLeft(2)
```

Use `style.Render("text")` in `View()` to apply styling.

---

## 2. Process Tree Walking on macOS

### ps Command

```bash
ps -ax -o pid,ppid,comm
```

Output includes a header line, then space-separated columns:
```
  PID  PPID COMM
    1     0 launchd
  256     1 sshd
  257   256 bash
  300   257 claude
```

- `ps -ax` and `ps -e` are equivalent on macOS (all processes)
- `comm` returns the **base executable name** (e.g., `claude`, `bash`), not the full path — ideal for matching
- Header line must be skipped when parsing

### Parsing and Tree Building in Go

```go
func buildProcessTree() (map[int][]int, map[int]string, error) {
    out, err := exec.Command("ps", "-ax", "-o", "pid,ppid,comm").Output()
    if err != nil {
        return nil, nil, err
    }

    children := make(map[int][]int)   // parent → children
    names := make(map[int]string)      // pid → command name

    scanner := bufio.NewScanner(bytes.NewReader(out))
    scanner.Scan() // skip header

    for scanner.Scan() {
        fields := strings.Fields(scanner.Text())
        if len(fields) < 3 { continue }

        pid, err1 := strconv.Atoi(fields[0])
        ppid, err2 := strconv.Atoi(fields[1])
        if err1 != nil || err2 != nil { continue }

        children[ppid] = append(children[ppid], pid)
        names[pid] = fields[2]
    }
    return children, names, nil
}
```

### Walking Descendants (BFS)

Given a pane's shell PID from tmux, walk all descendants looking for `claude`:

```go
func hasClaudeDescendant(rootPID int, children map[int][]int, names map[int]string) bool {
    queue := []int{rootPID}
    for len(queue) > 0 {
        pid := queue[0]
        queue = queue[1:]
        if names[pid] == "claude" {
            return true
        }
        queue = append(queue, children[pid]...)
    }
    return false
}
```

**Key design choice:** Build the full process tree once per refresh (single `ps` call), then walk it for each pane. Avoids N+1 subprocess calls.

### Edge Cases
- Zombie processes (`Z` state) still appear in ps with a command name — they won't cause false positives since they'd only match if named `claude`
- Process may exit between `ps` and use — not a problem since we only read from the snapshot
- macOS background processes often reparent to launchd (PID 1) — irrelevant since we walk down from known pane PIDs, not up

---

## 3. tmux CLI Output Parsing

### list-sessions

```bash
tmux list-sessions -F '#{session_name}|#{session_attached}|#{session_windows}'
```

Output (one line per session):
```
webapp|1|3
build|0|1
```

**Available format variables (confirmed in modern tmux):**
- `#{session_name}` — session name
- `#{session_attached}` — `1` if attached, `0` if not
- `#{session_windows}` — window count

**IMPORTANT:** `#{session_width}` and `#{session_height}` were **removed in tmux 2.9**. Use `#{window_width}` / `#{window_height}` or `#{pane_width}` / `#{pane_height}` instead.

### list-panes

```bash
tmux list-panes -t <session> -F '#{pane_pid}|#{pane_current_path}'
```

Output:
```
12345|/Users/alice/project
12400|/Users/alice/other
```

- `#{pane_pid}` returns the **shell PID** (the process launched in the pane — typically bash/zsh). Claude Code runs as a descendant of this PID.
- `#{pane_current_path}` returns the pane's current working directory.

### Delimiter Choice

**Use pipe `|` instead of colon `:`.** File paths contain colons on some systems, and `|` is never valid in paths. This avoids ambiguous parsing.

### Parsing in Go

```go
func listSessions() ([]Session, error) {
    out, err := exec.Command("tmux", "list-sessions",
        "-F", "#{session_name}|#{session_attached}|#{session_windows}").Output()
    if err != nil {
        return nil, err
    }

    var sessions []Session
    for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
        if line == "" { continue }
        parts := strings.SplitN(line, "|", 3)
        if len(parts) < 3 { continue }
        attached, _ := strconv.Atoi(parts[1])
        windows, _ := strconv.Atoi(parts[2])
        sessions = append(sessions, Session{
            Name:     parts[0],
            Attached: attached > 0,
            Windows:  windows,
        })
    }
    return sessions, nil
}
```

### Error Handling

| Scenario | Exit Code | Stderr |
|---|---|---|
| No tmux server running | 1 | `no server running on ...` |
| Server running, no sessions | 1 | `no sessions` |
| Session not found (list-panes) | 1 | `can't find session <name>` |
| Success | 0 | (empty) |

Check `err != nil` from `cmd.Output()` — both "no server" and "no sessions" return exit code 1. Parse stderr to distinguish if needed, or treat both as "no sessions available."

### attach-session

```bash
tmux attach-session -t <session>
```

Flags:
- `-d` — detach other clients first (force-attach, useful for taking over a session)
- `-r` — read-only mode

---

## 4. Cobra CLI Subcommand Patterns

### Project Layout

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
    ├── tmux/         # tmux interaction (list, attach, parse)
    ├── process/      # process tree walking
    ├── tui/          # Bubble Tea TUI
    └── profile/      # shell profile injection
```

### Root Command (Runs TUI)

```go
var rootCmd = &cobra.Command{
    Use:   "rta",
    Short: "Remote tmux access — find and attach to Claude Code sessions",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runTUI()  // Launch TUI when no subcommand given
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(attachCmd)
    rootCmd.AddCommand(statusCmd)
    rootCmd.AddCommand(setupCmd)
}
```

### Positional Argument (attach)

```go
var attachCmd = &cobra.Command{
    Use:   "attach <name>",
    Short: "Attach to a tmux session by name (substring match)",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return attachToSession(args[0])
    },
}
```

### Boolean Flags (--json, --undo)

```go
var jsonOutput bool

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "List sessions (one per line)",
    RunE: func(cmd *cobra.Command, args []string) error {
        return showStatus(jsonOutput)
    },
}

func init() {
    statusCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "output as JSON")
}
```

### Nested Subcommand (setup ssh)

```go
// setup.go
var undoSetup bool
var setupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Add auto-launch block to shell profile",
    RunE: func(cmd *cobra.Command, args []string) error {
        if undoSetup { return removeAutoLaunch() }
        return installAutoLaunch()
    },
}

func init() {
    setupCmd.Flags().BoolVar(&undoSetup, "undo", false, "remove auto-launch block")
    setupCmd.AddCommand(setupSSHCmd)
}

// setup_ssh.go
var setupSSHCmd = &cobra.Command{
    Use:   "ssh",
    Short: "Check SSH server configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        return checkSSHConfig()
    },
}
```

### main.go

```go
package main

import "rta/cmd"

func main() {
    cmd.Execute()
}
```

---

## 5. syscall.Exec for Process Replacement

### When to Use

- **`syscall.Exec`**: For `rta attach <name>` — replaces the Go process entirely with tmux. The user's shell prompt reappears after tmux detach. Cleanest UX.
- **`tea.ExecProcess`**: For the TUI case — suspends the TUI, runs tmux as a child, resumes TUI on detach.

### Implementation

```go
func attachSession(sessionName string) error {
    binary, err := exec.LookPath("tmux")
    if err != nil {
        return fmt.Errorf("tmux not found: %w", err)
    }

    args := []string{"tmux", "attach-session", "-t", sessionName}
    return syscall.Exec(binary, args, os.Environ())
}
```

**Critical details:**
- `exec.LookPath("tmux")` returns the full path (e.g., `/usr/bin/tmux`) — `syscall.Exec` requires an absolute path
- `args[0]` must be the program name (`"tmux"`) — Unix convention
- `os.Environ()` passes the current environment
- `syscall.Exec` **only returns on error** — on success, the Go process ceases to exist
- Works on macOS and Linux; not available on Windows (not needed)

### Why syscall.Exec over os/exec.Command

| | syscall.Exec | os/exec.Command |
|---|---|---|
| Process | Replaces current process | Spawns child process |
| After exit | Shell prompt returns | Go process resumes |
| PID | Same PID as rta | New PID |
| Use case | Final attach (no return) | Need to continue after |

---

## 6. Shell Profile Block Injection/Removal

### The Block

```bash
# rta: remote tmux access (auto-launch)
if [ -n "$SSH_CONNECTION" ] && command -v rta >/dev/null 2>&1; then
  rta
fi
# end rta
```

### Shell Detection

```go
func profilePath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil { return "", err }

    shell := os.Getenv("SHELL")
    switch {
    case strings.HasSuffix(shell, "zsh"):
        return filepath.Join(home, ".zshrc"), nil
    case strings.HasSuffix(shell, "bash"):
        return filepath.Join(home, ".bashrc"), nil
    default:
        return "", fmt.Errorf("unsupported shell: %s", shell)
    }
}
```

### Injection (Idempotent)

```go
const startMarker = "# rta: remote tmux access (auto-launch)"
const endMarker = "# end rta"
const block = `# rta: remote tmux access (auto-launch)
if [ -n "$SSH_CONNECTION" ] && command -v rta >/dev/null 2>&1; then
  rta
fi
# end rta`

func installAutoLaunch() error {
    path, err := profilePath()
    if err != nil { return err }

    content, err := os.ReadFile(path)
    if err != nil && !errors.Is(err, os.ErrNotExist) {
        return err
    }

    if strings.Contains(string(content), startMarker) {
        fmt.Println("Auto-launch block already installed.")
        return nil
    }

    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil { return err }
    defer f.Close()

    s := string(content)
    if len(s) > 0 && !strings.HasSuffix(s, "\n") {
        f.WriteString("\n")
    }
    _, err = f.WriteString("\n" + block + "\n")
    return err
}
```

### Removal

```go
func removeAutoLaunch() error {
    path, err := profilePath()
    if err != nil { return err }

    content, err := os.ReadFile(path)
    if err != nil { return err }

    lines := strings.Split(string(content), "\n")
    var result []string
    removing := false

    for _, line := range lines {
        if strings.Contains(line, startMarker) {
            removing = true
            continue
        }
        if strings.Contains(line, endMarker) {
            removing = false
            continue
        }
        if !removing {
            result = append(result, line)
        }
    }

    info, _ := os.Stat(path)
    perm := info.Mode().Perm()
    return os.WriteFile(path, []byte(strings.Join(result, "\n")), perm)
}
```

### Safety Considerations
- **Idempotent**: Check for `startMarker` before injecting; no-op if already present
- **File creation**: `os.O_CREATE` handles missing `.zshrc` / `.bashrc`
- **Permission preservation**: Read existing permissions with `os.Stat()`, apply with `os.WriteFile()`
- **No trailing blank lines**: The removal loop preserves surrounding content exactly

---

## Implementation Sequence

Recommended build order based on dependency analysis:

1. **tmux parsing** (`internal/tmux/`) — list sessions, list panes, parse output
2. **Process tree** (`internal/process/`) — build tree, walk descendants, detect claude
3. **Cobra commands** (`cmd/`) — wire up root, attach, status, setup
4. **syscall.Exec attach** — `rta attach <name>` working end-to-end
5. **Bubble Tea TUI** (`internal/tui/`) — session list with cursor, Enter to attach via tea.ExecProcess
6. **Shell profile injection** (`internal/profile/`) — setup / setup --undo
7. **SSH config checker** — setup ssh (inspect-only, no writes)

---

## Key Dependencies

```
github.com/charmbracelet/bubbletea  v2   — TUI framework
github.com/charmbracelet/lipgloss   v2   — TUI styling
github.com/spf13/cobra                   — CLI routing
```

No other external dependencies needed. All tmux/process interaction uses `os/exec` and standard library.
