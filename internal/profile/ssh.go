package profile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CheckSSHConfig() error {
	fmt.Println("SSH Configuration Check:")

	checkSSHD()
	checkAuthorizedKeys()
	checkSSHDirPerms()
	checkAuthorizedKeysPerms()

	return nil
}

func checkSSHD() {
	out, err := exec.Command("pgrep", "-x", "sshd").Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		fmt.Println("  [!!] sshd is not running")
		fmt.Println("       Fix: Enable Remote Login in System Settings > General > Sharing")
		return
	}
	fmt.Println("  [OK] sshd is running")
}

func checkAuthorizedKeys() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("  [!!] cannot determine home directory: %v\n", err)
		return
	}

	path := filepath.Join(home, ".ssh", "authorized_keys")
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("  [!!] ~/.ssh/authorized_keys does not exist")
		fmt.Println("       Fix: mkdir -p ~/.ssh && touch ~/.ssh/authorized_keys")
		return
	}
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			count++
		}
	}

	if count == 0 {
		fmt.Println("  [!!] ~/.ssh/authorized_keys exists but is empty")
		fmt.Println("       Fix: Add your public key to ~/.ssh/authorized_keys")
		return
	}
	fmt.Printf("  [OK] ~/.ssh/authorized_keys exists (%d keys)\n", count)
}

func checkSSHDirPerms() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	dir := filepath.Join(home, ".ssh")
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Println("  [!!] ~/.ssh directory does not exist")
		fmt.Println("       Fix: mkdir -p ~/.ssh && chmod 700 ~/.ssh")
		return
	}

	perm := info.Mode().Perm()
	if perm != 0700 {
		fmt.Printf("  [!!] ~/.ssh permissions are %04o (should be 0700)\n", perm)
		fmt.Println("       Fix: chmod 700 ~/.ssh")
		return
	}
	fmt.Println("  [OK] ~/.ssh permissions are 0700")
}

func checkAuthorizedKeysPerms() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	path := filepath.Join(home, ".ssh", "authorized_keys")
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		fmt.Printf("  [!!] ~/.ssh/authorized_keys permissions are %04o (should be 0600)\n", perm)
		fmt.Println("       Fix: chmod 600 ~/.ssh/authorized_keys")
		return
	}
	fmt.Println("  [OK] ~/.ssh/authorized_keys permissions are 0600")
}
