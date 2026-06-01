package cmd

import (
	"errors"
	"fmt"
	"log/slog"

	"rta/internal/tmux"

	"github.com/spf13/cobra"
)

var attachCmd = &cobra.Command{
	Use:   "attach <name>",
	Short: "Attach to a tmux session by name (substring match)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return attachToSession(args[0])
	},
}

func attachToSession(name string) error {
	sessions, err := tmux.ListSessions()
	if err != nil {
		if errors.Is(err, tmux.ErrNoServer) {
			return fmt.Errorf("no tmux sessions found (tmux server is not running)")
		}
		if errors.Is(err, tmux.ErrNoSessions) {
			return fmt.Errorf("no tmux sessions found")
		}
		return fmt.Errorf("listing sessions: %w", err)
	}
	if len(sessions) == 0 {
		return fmt.Errorf("no tmux sessions found")
	}

	matches := tmux.FindSession(name, sessions)
	slog.Debug("searching for session", "query", name, "match_count", len(matches))

	switch len(matches) {
	case 0:
		fmt.Printf("No sessions matching %q\n", name)
		fmt.Println("\nAvailable sessions:")
		for _, s := range sessions {
			fmt.Printf("  %s\n", s.Name)
		}
		return fmt.Errorf("no matching session")
	case 1:
		if err := tmux.AttachSession(matches[0].Name); err != nil {
			return fmt.Errorf("attaching to session %q: %w", matches[0].Name, err)
		}
		return nil
	default:
		fmt.Printf("Multiple sessions match %q — be more specific:\n", name)
		for _, s := range matches {
			fmt.Printf("  %s\n", s.Name)
		}
		return fmt.Errorf("ambiguous session name")
	}
}
