package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var undoSetup bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Add auto-launch block to shell profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("setup: not yet implemented")
		return nil
	},
}

func init() {
	setupCmd.Flags().BoolVar(&undoSetup, "undo", false, "remove auto-launch block")
	setupCmd.AddCommand(setupSSHCmd)
}
