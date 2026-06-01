package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"rta/internal/process"
	"rta/internal/tmux"
)

func resetRootCmd(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		verbose = false
	})
}

func TestStatusCmdNoSessions(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoServer })

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"status"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("status: %v", err)
	}
}

func TestStatusCmdJSON(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoSessions })

	rootCmd.SetArgs([]string{"status", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("status --json: %v", err)
	}
}

func TestStatusCmdWithSessions(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|1|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) { return []byte("100|/home/user\n"), nil })
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n  101   100 claude\n"), nil
	})

	rootCmd.SetArgs([]string{"status"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("status: %v", err)
	}
}

func TestStatusCmdWithSessionsJSON(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) { return []byte("100|/home/user\n"), nil })
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n  101   100 claude\n"), nil
	})

	rootCmd.SetArgs([]string{"status", "-j"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("status -j: %v", err)
	}
}

func TestAttachCmdNoArgs(t *testing.T) {
	resetRootCmd(t)
	rootCmd.SetArgs([]string{"attach"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach without args should fail")
	}
}

func TestAttachCmdNoServer(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoServer })

	rootCmd.SetArgs([]string{"attach", "test"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail when no server")
	}
}

func TestAttachCmdNoSessions(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoSessions })

	rootCmd.SetArgs([]string{"attach", "test"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail when no sessions")
	}
}

func TestAttachCmdEmptySessions(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte(""), nil })

	rootCmd.SetArgs([]string{"attach", "test"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail when empty sessions")
	}
}

func TestAttachCmdNoMatch(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })

	rootCmd.SetArgs([]string{"attach", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail when no match")
	}
}

func TestAttachCmdAmbiguous(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("webapp|0|1\nweb-api|0|1\n"), nil })

	rootCmd.SetArgs([]string{"attach", "web"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail when ambiguous")
	}
}

func TestAttachCmdMatch(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetLookPath(t, func(file string) (string, error) { return "/usr/bin/tmux", nil })
	tmux.SetExecSyscall(t, func(binary string, args []string, env []string) error { return nil })

	rootCmd.SetArgs([]string{"attach", "dev"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("attach should succeed: %v", err)
	}
}

func TestAttachCmdAttachError(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetLookPath(t, func(file string) (string, error) { return "/usr/bin/tmux", nil })
	tmux.SetExecSyscall(t, func(binary string, args []string, env []string) error { return fmt.Errorf("attach failed") })

	rootCmd.SetArgs([]string{"attach", "dev"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail when exec fails")
	}
}

func TestAttachCmdListError(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, fmt.Errorf("tmux: some other error") })

	rootCmd.SetArgs([]string{"attach", "test"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("attach should fail on list error")
	}
}

func TestStatusCmdProcessTreeError(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	process.SetRunPS(t, func() ([]byte, error) { return nil, fmt.Errorf("ps failed") })

	rootCmd.SetArgs([]string{"status"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("status should fail when process tree fails")
	}
}

func TestSingleClaudeSessionNone(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoServer })

	if _, ok := singleClaudeSession(); ok {
		t.Error("expected false when no server")
	}
}

func TestSingleClaudeSessionNoSessions(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte(""), nil })

	if _, ok := singleClaudeSession(); ok {
		t.Error("expected false when no sessions")
	}
}

func TestSingleClaudeSessionBuildTreeError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	process.SetRunPS(t, func() ([]byte, error) { return nil, fmt.Errorf("ps failed") })

	if _, ok := singleClaudeSession(); ok {
		t.Error("expected false when build tree fails")
	}
}

func TestSingleClaudeSessionOne(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\nbuild|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) {
		if name == "dev" {
			return []byte("100|/home\n"), nil
		}
		return []byte("200|/home\n"), nil
	})
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n  101   100 claude\n  200     1 bash\n"), nil
	})

	sess, ok := singleClaudeSession()
	if !ok || sess != "dev" {
		t.Errorf("expected (dev, true), got (%q, %v)", sess, ok)
	}
}

func TestSingleClaudeSessionMultiple(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\nbuild|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) {
		if name == "dev" {
			return []byte("100|/home\n"), nil
		}
		return []byte("200|/home\n"), nil
	})
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n  101   100 claude\n  200     1 bash\n  201   200 claude\n"), nil
	})

	if _, ok := singleClaudeSession(); ok {
		t.Error("expected false for multiple claude sessions")
	}
}

func TestSingleClaudeSessionPaneError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) { return nil, fmt.Errorf("pane error") })
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n"), nil
	})

	if _, ok := singleClaudeSession(); ok {
		t.Error("expected false when pane error")
	}
}

func TestSingleClaudeSessionNoClaude(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) { return []byte("100|/home\n"), nil })
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n"), nil
	})

	if _, ok := singleClaudeSession(); ok {
		t.Error("expected false when no claude")
	}
}

func TestVerboseFlag(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoServer })

	rootCmd.SetArgs([]string{"status", "-v"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("status -v: %v", err)
	}
}

func TestShowStatusEmptySessions(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte(""), nil })

	if err := showStatus(false); err != nil {
		t.Fatalf("showStatus: %v", err)
	}
}

func TestShowStatusEmptySessionsJSON(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte(""), nil })

	if err := showStatus(true); err != nil {
		t.Fatalf("showStatus json: %v", err)
	}
}

func TestShowStatusGenericError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, fmt.Errorf("tmux: unknown error") })

	err := showStatus(false)
	if err == nil || !strings.Contains(err.Error(), "listing sessions") {
		t.Errorf("expected 'listing sessions' error, got: %v", err)
	}
}

func TestShowStatusNoServerJSON(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoServer })

	if err := showStatus(true); err != nil {
		t.Fatalf("showStatus: %v", err)
	}
}

func TestShowStatusNoServerText(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoSessions })

	if err := showStatus(false); err != nil {
		t.Fatalf("showStatus: %v", err)
	}
}

func TestShowStatusPaneError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) { return nil, fmt.Errorf("pane error") })
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n"), nil
	})

	if err := showStatus(false); err != nil {
		t.Fatalf("showStatus should handle pane error: %v", err)
	}
}

func TestSetupCmdInstall(t *testing.T) {
	resetRootCmd(t)
	t.Setenv("SHELL", "/bin/zsh")
	t.Setenv("HOME", t.TempDir())
	rootCmd.SetArgs([]string{"setup"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("setup: %v", err)
	}
}

func TestSetupCmdUndo(t *testing.T) {
	resetRootCmd(t)
	t.Setenv("SHELL", "/bin/bash")
	t.Setenv("HOME", t.TempDir())
	rootCmd.SetArgs([]string{"setup", "--undo"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("setup --undo: %v", err)
	}
}

func TestSetupSSHCmd(t *testing.T) {
	resetRootCmd(t)
	rootCmd.SetArgs([]string{"setup", "ssh"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("setup ssh: %v", err)
	}
}

func TestExecuteFunction(t *testing.T) {
	resetRootCmd(t)
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, tmux.ErrNoServer })

	rootCmd.SetArgs([]string{"status"})
	Execute()
}
