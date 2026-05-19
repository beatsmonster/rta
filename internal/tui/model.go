package tui

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SessionInfo struct {
	Name       string
	HasClaude  bool
	WorkingDir string
	Attached   bool
}

type sessionsMsg struct {
	sessions []SessionInfo
	err      error
}

type tmuxFinishedMsg struct{ err error }

type model struct {
	sessions []SessionInfo
	cursor   int
	width    int
	height   int
	err      error
}

func New() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return refreshSessions
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.sessions) > 0 {
				return m, attachTmux(m.sessions[m.cursor].Name)
			}
		case "r":
			return m, refreshSessions
		case "q", "ctrl+c", "escape":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case sessionsMsg:
		if msg.err != nil {
			m.err = msg.err
			m.sessions = nil
			return m, nil
		}
		m.err = nil
		m.sessions = msg.sessions
		if m.cursor >= len(m.sessions) {
			m.cursor = max(0, len(m.sessions)-1)
		}
	case tmuxFinishedMsg:
		return m, refreshSessions
	}
	return m, nil
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	normalStyle   = lipgloss.NewStyle()
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView(errStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\n" + helpStyle.Render("[r] retry  [esc/q] quit"))
	}

	if len(m.sessions) == 0 {
		return tea.NewView("No tmux sessions found.\n\n" + helpStyle.Render("[r] refresh  [esc/q] quit"))
	}

	var claudeSessions, otherSessions []int
	for i, s := range m.sessions {
		if s.HasClaude {
			claudeSessions = append(claudeSessions, i)
		} else {
			otherSessions = append(otherSessions, i)
		}
	}

	var b strings.Builder

	if len(claudeSessions) > 0 {
		b.WriteString(titleStyle.Render("Claude Code Sessions"))
		b.WriteString("\n")
		for _, i := range claudeSessions {
			b.WriteString(m.renderSession(i))
			b.WriteString("\n")
		}
	}

	if len(otherSessions) > 0 {
		if len(claudeSessions) > 0 {
			b.WriteString("\n")
		}
		b.WriteString(titleStyle.Render("Other tmux Sessions"))
		b.WriteString("\n")
		for _, i := range otherSessions {
			b.WriteString(m.renderSession(i))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[enter] attach  [r] refresh  [esc/q] quit"))

	return tea.NewView(b.String())
}

func (m model) renderSession(idx int) string {
	s := m.sessions[idx]
	cursor := "  "
	if idx == m.cursor {
		cursor = "> "
	}

	dir := s.WorkingDir
	if home, err := os.UserHomeDir(); err == nil {
		dir = strings.Replace(dir, home, "~", 1)
	}

	line := fmt.Sprintf("%s%-16s %s", cursor, s.Name, dir)
	if s.Attached {
		line += " (attached)"
	}

	if idx == m.cursor {
		return selectedStyle.Render(line)
	}
	return normalStyle.Render(line)
}
