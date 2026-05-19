package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestInstallSuccess(t *testing.T) {
	t.Setenv("SHELL", "/bin/zsh")
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.WriteFile(filepath.Join(dir, ".zshrc"), []byte("# existing\n"), 0644)

	if err := Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	content, _ := os.ReadFile(filepath.Join(dir, ".zshrc"))
	if !strings.Contains(string(content), startMarker) {
		t.Error("Install did not add block")
	}
}

func TestRemoveSuccess(t *testing.T) {
	t.Setenv("SHELL", "/bin/bash")
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.WriteFile(filepath.Join(dir, ".bashrc"), []byte("before\n"+block+"\nafter\n"), 0644)

	if err := Remove(); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	content, _ := os.ReadFile(filepath.Join(dir, ".bashrc"))
	if strings.Contains(string(content), startMarker) {
		t.Error("Remove did not remove block")
	}
}

func TestInstallUnsupportedShell(t *testing.T) {
	t.Setenv("SHELL", "/bin/fish")
	if err := Install(); err == nil {
		t.Error("expected error for unsupported shell")
	}
}

func TestRemoveUnsupportedShell(t *testing.T) {
	t.Setenv("SHELL", "/bin/fish")
	if err := Remove(); err == nil {
		t.Error("expected error for unsupported shell")
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	return string(buf[:n])
}

func TestCheckSSHConfigRuns(t *testing.T) {
	dir := t.TempDir()
	sshDir := filepath.Join(dir, ".ssh")
	os.MkdirAll(sshDir, 0700)
	os.WriteFile(filepath.Join(sshDir, "authorized_keys"), []byte("ssh-rsa AAAA...\n"), 0600)

	SetRunPgrep(t, func() ([]byte, error) { return []byte("1234\n"), nil })
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() {
		if err := CheckSSHConfig(); err != nil {
			t.Fatalf("CheckSSHConfig: %v", err)
		}
	})
	if !strings.Contains(out, "[OK] sshd is running") {
		t.Errorf("expected sshd OK, got: %s", out)
	}
}

func TestCheckSSHDNotRunning(t *testing.T) {
	SetRunPgrep(t, func() ([]byte, error) { return nil, fmt.Errorf("no sshd") })

	out := captureStdout(t, func() { checkSSHD() })
	if !strings.Contains(out, "[!!] sshd is not running") {
		t.Errorf("expected not running, got: %s", out)
	}
}

func TestCheckSSHDRunning(t *testing.T) {
	SetRunPgrep(t, func() ([]byte, error) { return []byte("1234\n"), nil })

	out := captureStdout(t, func() { checkSSHD() })
	if !strings.Contains(out, "[OK] sshd is running") {
		t.Errorf("expected running, got: %s", out)
	}
}

func TestCheckSSHDEmptyOutput(t *testing.T) {
	SetRunPgrep(t, func() ([]byte, error) { return []byte("  \n"), nil })

	out := captureStdout(t, func() { checkSSHD() })
	if !strings.Contains(out, "[!!] sshd is not running") {
		t.Errorf("expected not running for empty, got: %s", out)
	}
}

func TestCheckAuthorizedKeysHomeDirError(t *testing.T) {
	SetGetHomeDir(t, func() (string, error) { return "", fmt.Errorf("no home") })

	out := captureStdout(t, func() { checkAuthorizedKeys() })
	if !strings.Contains(out, "cannot determine home directory") {
		t.Errorf("expected home dir error, got: %s", out)
	}
}

func TestCheckAuthorizedKeysNotExists(t *testing.T) {
	SetGetHomeDir(t, func() (string, error) { return t.TempDir(), nil })

	out := captureStdout(t, func() { checkAuthorizedKeys() })
	if !strings.Contains(out, "does not exist") {
		t.Errorf("expected not exist, got: %s", out)
	}
}

func TestCheckAuthorizedKeysEmpty(t *testing.T) {
	dir := t.TempDir()
	sshDir := filepath.Join(dir, ".ssh")
	os.MkdirAll(sshDir, 0700)
	os.WriteFile(filepath.Join(sshDir, "authorized_keys"), []byte("# comment\n\n"), 0600)
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() { checkAuthorizedKeys() })
	if !strings.Contains(out, "exists but is empty") {
		t.Errorf("expected empty, got: %s", out)
	}
}

func TestCheckAuthorizedKeysWithKeys(t *testing.T) {
	dir := t.TempDir()
	sshDir := filepath.Join(dir, ".ssh")
	os.MkdirAll(sshDir, 0700)
	os.WriteFile(filepath.Join(sshDir, "authorized_keys"), []byte("ssh-rsa AAAA key1\nssh-ed25519 AAAB key2\n"), 0600)
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() { checkAuthorizedKeys() })
	if !strings.Contains(out, "(2 keys)") {
		t.Errorf("expected 2 keys, got: %s", out)
	}
}

func TestCheckSSHDirPermsHomeDirError(t *testing.T) {
	SetGetHomeDir(t, func() (string, error) { return "", fmt.Errorf("no home") })

	out := captureStdout(t, func() { checkSSHDirPerms() })
	if out != "" {
		t.Errorf("expected no output, got: %s", out)
	}
}

func TestCheckSSHDirPermsNotExists(t *testing.T) {
	SetGetHomeDir(t, func() (string, error) { return t.TempDir(), nil })

	out := captureStdout(t, func() { checkSSHDirPerms() })
	if !strings.Contains(out, "does not exist") {
		t.Errorf("expected not exist, got: %s", out)
	}
}

func TestCheckSSHDirPermsWrong(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".ssh"), 0755)
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() { checkSSHDirPerms() })
	if !strings.Contains(out, "should be 0700") {
		t.Errorf("expected wrong perms, got: %s", out)
	}
}

func TestCheckSSHDirPermsCorrect(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".ssh"), 0700)
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() { checkSSHDirPerms() })
	if !strings.Contains(out, "[OK] ~/.ssh permissions are 0700") {
		t.Errorf("expected OK, got: %s", out)
	}
}

func TestCheckAuthorizedKeysPermsHomeDirError(t *testing.T) {
	SetGetHomeDir(t, func() (string, error) { return "", fmt.Errorf("no home") })

	out := captureStdout(t, func() { checkAuthorizedKeysPerms() })
	if out != "" {
		t.Errorf("expected no output, got: %s", out)
	}
}

func TestCheckAuthorizedKeysPermsNotExists(t *testing.T) {
	SetGetHomeDir(t, func() (string, error) { return t.TempDir(), nil })

	out := captureStdout(t, func() { checkAuthorizedKeysPerms() })
	if out != "" {
		t.Errorf("expected no output for missing file, got: %s", out)
	}
}

func TestCheckAuthorizedKeysPermsWrong(t *testing.T) {
	dir := t.TempDir()
	sshDir := filepath.Join(dir, ".ssh")
	os.MkdirAll(sshDir, 0700)
	os.WriteFile(filepath.Join(sshDir, "authorized_keys"), []byte("key\n"), 0644)
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() { checkAuthorizedKeysPerms() })
	if !strings.Contains(out, "should be 0600") {
		t.Errorf("expected wrong perms, got: %s", out)
	}
}

func TestCheckAuthorizedKeysPermsCorrect(t *testing.T) {
	dir := t.TempDir()
	sshDir := filepath.Join(dir, ".ssh")
	os.MkdirAll(sshDir, 0700)
	os.WriteFile(filepath.Join(sshDir, "authorized_keys"), []byte("key\n"), 0600)
	SetGetHomeDir(t, func() (string, error) { return dir, nil })

	out := captureStdout(t, func() { checkAuthorizedKeysPerms() })
	if !strings.Contains(out, "[OK] ~/.ssh/authorized_keys permissions are 0600") {
		t.Errorf("expected OK, got: %s", out)
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
