package cmd

import (
	"fmt"

	"rta/internal/profile"

	"github.com/spf13/cobra"
)

var undoSetup bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Add auto-launch block to shell profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		if undoSetup {
			if err := profile.Remove(); err != nil {
				return err
			}
			fmt.Println("Auto-launch block removed.")
			return nil
		}
		if err := profile.Install(); err != nil {
			return err
		}
		fmt.Println("Auto-launch block installed.")
		return nil
	},
}

func init() {
	setupCmd.Flags().BoolVar(&undoSetup, "undo", false, "remove auto-launch block")
	setupCmd.AddCommand(setupSSHCmd)
}
