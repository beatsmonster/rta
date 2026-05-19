package tui

import (
	"log/slog"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"rta/internal/process"
	"rta/internal/recap"
	"rta/internal/tmux"
)

func refreshSessions() tea.Msg {
	sessions, err := tmux.ListSessions()
	if err != nil {
		return sessionsMsg{err: err}
	}

	children, names, err := process.BuildTree()
	if err != nil {
		return sessionsMsg{err: err}
	}

	var infos []SessionInfo
	for _, s := range sessions {
		info := SessionInfo{
			Name:     s.Name,
			Attached: s.Attached,
		}

		panes, err := tmux.ListPanes(s.Name)
		if err == nil && len(panes) > 0 {
			info.WorkingDir = panes[0].CurrentPath
			for _, p := range panes {
				if process.HasDescendant(p.PID, "claude", children, names) {
					info.HasClaude = true
					break
				}
			}
		}

		infos = append(infos, info)
	}

	slog.Debug("refreshed sessions", "session_count", len(infos))
	return sessionsMsg{sessions: infos}
}

func loadRecaps() tea.Msg {
	recaps, err := recap.List()
	if err != nil {
		slog.Debug("failed to load recaps", "error", err)
		return recapsMsg{}
	}
	m := make(map[string]*recap.Recap, len(recaps))
	for _, r := range recaps {
		m[r.Name] = r
	}
	return recapsMsg{recaps: m}
}

func attachTmux(name string) tea.Cmd {
	c := exec.Command("tmux", "attach-session", "-t", name)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return tmuxFinishedMsg{err}
	})
}
