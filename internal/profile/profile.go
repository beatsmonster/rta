package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const startMarker = "# rta: remote tmux access (auto-launch)"
const endMarker = "# end rta"
const block = `# rta: remote tmux access (auto-launch)
if [ -n "$SSH_CONNECTION" ] && command -v rta >/dev/null 2>&1; then
  rta
fi
# end rta`

func ProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	shell := os.Getenv("SHELL")
	switch {
	case strings.HasSuffix(shell, "zsh"):
		return filepath.Join(home, ".zshrc"), nil
	case strings.HasSuffix(shell, "bash"):
		return filepath.Join(home, ".bashrc"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

func Install() error {
	path, err := ProfilePath()
	if err != nil {
		return err
	}
	return InstallTo(path)
}

func Remove() error {
	path, err := ProfilePath()
	if err != nil {
		return err
	}
	return RemoveFrom(path)
}

func InstallTo(path string) error {
	content, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if strings.Contains(string(content), startMarker) {
		return nil
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	s := string(content)
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		f.WriteString("\n")
	}
	_, err = f.WriteString("\n" + block + "\n")
	return err
}

func RemoveFrom(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if !strings.Contains(string(content), startMarker) {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	perm := info.Mode().Perm()

	lines := strings.Split(string(content), "\n")
	var result []string
	removing := false

	for _, line := range lines {
		if strings.Contains(line, startMarker) {
			removing = true
			continue
		}
		if strings.Contains(line, endMarker) {
			removing = false
			continue
		}
		if !removing {
			result = append(result, line)
		}
	}

	return os.WriteFile(path, []byte(strings.Join(result, "\n")), perm)
}
