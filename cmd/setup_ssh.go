package cmd

import (
	"rta/internal/profile"

	"github.com/spf13/cobra"
)

var setupSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Check SSH server configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return profile.CheckSSHConfig()
	},
}
