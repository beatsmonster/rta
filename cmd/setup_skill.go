package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed skill_embed/SKILL.md
var skillContent []byte

var setupSkillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Install /enable-rta skill for Claude Code",
	RunE:  runSetupSkill,
}

func init() {
	setupCmd.AddCommand(setupSkillCmd)
}

func runSetupSkill(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}

	dir := filepath.Join(home, ".claude", "skills", "enable-rta")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating skill directory: %w", err)
	}

	dest := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(dest, skillContent, 0o644); err != nil {
		return fmt.Errorf("writing skill file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Installed /enable-rta skill to %s\n", dest)
	fmt.Fprintln(cmd.OutOrStdout(), "Usage: type /enable-rta in Claude Code to move your session into tmux.")
	return nil
}
