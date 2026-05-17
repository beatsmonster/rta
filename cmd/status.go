package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"rta/internal/process"
	"rta/internal/tmux"

	"github.com/spf13/cobra"
)

var jsonOutput bool

type sessionStatus struct {
	Name       string `json:"name"`
	HasClaude  bool   `json:"hasClaude"`
	WorkingDir string `json:"workingDir"`
	Attached   bool   `json:"attached"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List tmux sessions with Claude detection",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showStatus(jsonOutput)
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "output as JSON")
}

func showStatus(asJSON bool) error {
	sessions, err := tmux.ListSessions()
	if err != nil {
		return fmt.Errorf("listing sessions: %w", err)
	}

	if len(sessions) == 0 {
		if asJSON {
			fmt.Println("[]")
			return nil
		}
		fmt.Println("No tmux sessions found.")
		return nil
	}

	children, names, err := process.BuildTree()
	if err != nil {
		return fmt.Errorf("building process tree: %w", err)
	}

	var statuses []sessionStatus
	for _, s := range sessions {
		ss := sessionStatus{
			Name:     s.Name,
			Attached: s.Attached,
		}

		panes, err := tmux.ListPanes(s.Name)
		if err == nil && len(panes) > 0 {
			ss.WorkingDir = panes[0].CurrentPath
			for _, p := range panes {
				if process.HasDescendant(p.PID, "claude", children, names) {
					ss.HasClaude = true
					break
				}
			}
		}

		statuses = append(statuses, ss)
	}

	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(statuses)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "SESSION\tCLAUDE\tDIR\tATTACHED")
	for _, ss := range statuses {
		claude := " "
		if ss.HasClaude {
			claude = "*"
		}
		attached := " "
		if ss.Attached {
			attached = "yes"
		}
		dir := ss.WorkingDir
		if home, err := os.UserHomeDir(); err == nil {
			dir = strings.Replace(dir, home, "~", 1)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ss.Name, claude, dir, attached)
	}
	return w.Flush()
}
