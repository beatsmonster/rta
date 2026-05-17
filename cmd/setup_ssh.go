package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Check SSH server configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("setup ssh: not yet implemented")
		return nil
	},
}
