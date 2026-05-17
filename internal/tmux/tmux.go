package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Session struct {
	Name     string
	Attached bool
	Windows  int
}

type Pane struct {
	PID         int
	CurrentPath string
}

var ErrNoServer = fmt.Errorf("tmux: no server running")
var ErrNoSessions = fmt.Errorf("tmux: no sessions")

func ListSessions() ([]Session, error) {
	cmd := exec.Command("tmux", "list-sessions",
		"-F", "#{session_name}|#{session_attached}|#{session_windows}")
	out, err := cmd.Output()
	if err != nil {
		return nil, classifyError(err)
	}
	return parseSessions(string(out))
}

func ListPanes(sessionName string) ([]Pane, error) {
	cmd := exec.Command("tmux", "list-panes",
		"-t", sessionName,
		"-F", "#{pane_pid}|#{pane_current_path}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("tmux list-panes: %w", err)
	}
	return parsePanes(string(out))
}

func AttachSession(sessionName string) error {
	binary, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("tmux not found: %w", err)
	}
	args := []string{"tmux", "attach-session", "-t", sessionName}
	return syscall.Exec(binary, args, os.Environ())
}

func classifyError(err error) error {
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		stderr := string(exitErr.Stderr)
		if strings.Contains(stderr, "no server running") {
			return ErrNoServer
		}
		if strings.Contains(stderr, "no sessions") {
			return ErrNoSessions
		}
	}
	return fmt.Errorf("tmux: %w", err)
}

func parseSessions(output string) ([]Session, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	var sessions []Session
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}
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

func parsePanes(output string) ([]Pane, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	var panes []Pane
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) < 2 {
			continue
		}
		pid, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		panes = append(panes, Pane{
			PID:         pid,
			CurrentPath: parts[1],
		})
	}
	return panes, nil
}
