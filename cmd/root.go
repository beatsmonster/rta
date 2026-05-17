package cmd

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"rta/internal/process"
	"rta/internal/tmux"
	"rta/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "rta",
	Short: "Remote tmux Access — discover and attach to tmux sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		if sess, ok := singleClaudeSession(); ok {
			return tmux.AttachSession(sess)
		}
		_, err := tea.NewProgram(tui.New()).Run()
		return err
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

func singleClaudeSession() (string, bool) {
	sessions, err := tmux.ListSessions()
	if err != nil || len(sessions) == 0 {
		return "", false
	}

	children, names, err := process.BuildTree()
	if err != nil {
		return "", false
	}

	var claudeSession string
	claudeCount := 0

	for _, s := range sessions {
		panes, err := tmux.ListPanes(s.Name)
		if err != nil || len(panes) == 0 {
			continue
		}
		for _, p := range panes {
			if process.HasDescendant(p.PID, "claude", children, names) {
				claudeSession = s.Name
				claudeCount++
				break
			}
		}
		if claudeCount > 1 {
			return "", false
		}
	}

	if claudeCount == 1 {
		return claudeSession, true
	}
	return "", false
}
