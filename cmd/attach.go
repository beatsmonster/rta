package cmd

import (
	"fmt"
	"strings"

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
		return fmt.Errorf("listing sessions: %w", err)
	}

	var matches []tmux.Session
	for _, s := range sessions {
		if strings.Contains(s.Name, name) {
			matches = append(matches, s)
		}
	}

	switch len(matches) {
	case 0:
		fmt.Println("No sessions matching:", name)
		if len(sessions) > 0 {
			fmt.Println("\nAvailable sessions:")
			for _, s := range sessions {
				fmt.Printf("  %s\n", s.Name)
			}
		}
		return fmt.Errorf("no matching session")
	case 1:
		return tmux.AttachSession(matches[0].Name)
	default:
		fmt.Println("Multiple sessions match:", name)
		for _, s := range matches {
			fmt.Printf("  %s\n", s.Name)
		}
		return fmt.Errorf("ambiguous session name")
	}
}
