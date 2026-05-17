package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallToEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")

	if err := InstallTo(path); err != nil {
		t.Fatalf("InstallTo failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	s := string(content)
	if got := containsBlock(s); !got {
		t.Errorf("expected block in file, got:\n%s", s)
	}
}

func TestInstallToCreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".bashrc")

	if err := InstallTo(path); err != nil {
		t.Fatalf("InstallTo failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestInstallToIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")

	if err := InstallTo(path); err != nil {
		t.Fatalf("first InstallTo failed: %v", err)
	}
	first, _ := os.ReadFile(path)

	if err := InstallTo(path); err != nil {
		t.Fatalf("second InstallTo failed: %v", err)
	}
	second, _ := os.ReadFile(path)

	if string(first) != string(second) {
		t.Errorf("file changed on second install:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestInstallToExistingContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	existing := "export PATH=/usr/local/bin:$PATH\n"
	os.WriteFile(path, []byte(existing), 0644)

	if err := InstallTo(path); err != nil {
		t.Fatalf("InstallTo failed: %v", err)
	}

	content, _ := os.ReadFile(path)
	s := string(content)

	if !containsBlock(s) {
		t.Errorf("expected block in file")
	}
	if s[:len(existing)] != existing {
		t.Errorf("existing content was modified")
	}
}

func TestInstallToExistingContentNoTrailingNewline(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	existing := "export PATH=/usr/local/bin:$PATH"
	os.WriteFile(path, []byte(existing), 0644)

	if err := InstallTo(path); err != nil {
		t.Fatalf("InstallTo failed: %v", err)
	}

	content, _ := os.ReadFile(path)
	s := string(content)

	if !containsBlock(s) {
		t.Errorf("expected block in file")
	}
	if s[:len(existing)] != existing {
		t.Errorf("existing content was modified")
	}
}

func TestRemoveFromCleansBlock(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	before := "export FOO=bar\n"
	after := "export BAZ=qux\n"
	content := before + "\n" + block + "\n" + after
	os.WriteFile(path, []byte(content), 0644)

	if err := RemoveFrom(path); err != nil {
		t.Fatalf("RemoveFrom failed: %v", err)
	}

	result, _ := os.ReadFile(path)
	s := string(result)

	if containsBlock(s) {
		t.Errorf("block still present after removal:\n%s", s)
	}
	if !contains(s, "export FOO=bar") {
		t.Errorf("surrounding content was lost")
	}
	if !contains(s, "export BAZ=qux") {
		t.Errorf("surrounding content was lost")
	}
}

func TestRemoveFromNoBlock(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	content := "export FOO=bar\n"
	os.WriteFile(path, []byte(content), 0644)

	if err := RemoveFrom(path); err != nil {
		t.Fatalf("RemoveFrom failed: %v", err)
	}

	result, _ := os.ReadFile(path)
	if string(result) != content {
		t.Errorf("file was modified when no block present:\n%s", result)
	}
}

func TestRemoveFromNonexistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent")

	if err := RemoveFrom(path); err != nil {
		t.Fatalf("RemoveFrom should be no-op for missing file: %v", err)
	}
}

func TestRemoveFromPreservesPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	content := "before\n" + block + "\nafter\n"
	os.WriteFile(path, []byte(content), 0600)

	if err := RemoveFrom(path); err != nil {
		t.Fatalf("RemoveFrom failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("permissions changed: got %o, want 0600", info.Mode().Perm())
	}
}

func TestProfilePath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		shell string
		want  string
	}{
		{"/bin/zsh", filepath.Join(home, ".zshrc")},
		{"/usr/bin/zsh", filepath.Join(home, ".zshrc")},
		{"/bin/bash", filepath.Join(home, ".bashrc")},
		{"/usr/bin/bash", filepath.Join(home, ".bashrc")},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			t.Setenv("SHELL", tt.shell)
			got, err := ProfilePath()
			if err != nil {
				t.Fatalf("ProfilePath() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ProfilePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProfilePathUnsupported(t *testing.T) {
	t.Setenv("SHELL", "/bin/fish")
	_, err := ProfilePath()
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
}

func containsBlock(s string) bool {
	return contains(s, startMarker) && contains(s, endMarker)
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
