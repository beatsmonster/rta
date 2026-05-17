package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rta",
	Short: "Remote tmux Access — discover and attach to tmux sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("rta v0.1.0")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
