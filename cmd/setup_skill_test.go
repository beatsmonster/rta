package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkillContentEmbedded(t *testing.T) {
	if len(skillContent) == 0 {
		t.Fatal("embedded skill content is empty")
	}
	if !bytes.Contains(skillContent, []byte("name: enable-rta")) {
		t.Error("skill content missing expected frontmatter")
	}
}

func TestSetupSkillCreatesFile(t *testing.T) {
	resetRootCmd(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"setup", "skill"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("setup skill: %v", err)
	}

	dest := filepath.Join(home, ".claude", "skills", "enable-rta", "SKILL.md")
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("reading installed skill: %v", err)
	}
	if !bytes.Equal(data, skillContent) {
		t.Error("installed file content does not match embedded content")
	}
	if !strings.Contains(buf.String(), "Installed") {
		t.Error("expected success message in output")
	}
}

func TestSetupSkillOverwrite(t *testing.T) {
	resetRootCmd(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	dest := filepath.Join(home, ".claude", "skills", "enable-rta", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("old content"), 0o644); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"setup", "skill"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("setup skill overwrite: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "old content" {
		t.Error("skill file was not overwritten")
	}
	if !bytes.Equal(data, skillContent) {
		t.Error("overwritten content does not match embedded content")
	}
}
