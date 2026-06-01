package tui

import (
	"fmt"
	"testing"

	"rta/internal/process"
	"rta/internal/recap"
	"rta/internal/tmux"
)

func TestAttachTmux(t *testing.T) {
	if cmd := attachTmux("test"); cmd == nil {
		t.Fatal("attachTmux should return a non-nil cmd")
	}
}

func TestRefreshSessionsListError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return nil, fmt.Errorf("no server") })

	sm := refreshSessions().(sessionsMsg)
	if sm.err == nil {
		t.Error("expected error")
	}
}

func TestRefreshSessionsBuildTreeError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	process.SetRunPS(t, func() ([]byte, error) { return nil, fmt.Errorf("ps failed") })

	sm := refreshSessions().(sessionsMsg)
	if sm.err == nil {
		t.Error("expected error from BuildTree")
	}
}

func TestRefreshSessionsSuccess(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|1|1\nbuild|0|2\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) {
		if name == "dev" {
			return []byte("100|/home/user/project\n"), nil
		}
		return []byte("200|/home/user/build\n"), nil
	})
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n  101   100 claude\n  200     1 bash\n"), nil
	})

	sm := refreshSessions().(sessionsMsg)
	if sm.err != nil {
		t.Fatalf("unexpected error: %v", sm.err)
	}
	if len(sm.sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sm.sessions))
	}
	if !sm.sessions[0].HasClaude {
		t.Error("dev should have claude")
	}
	if sm.sessions[1].HasClaude {
		t.Error("build should not have claude")
	}
}

func TestLoadRecaps(t *testing.T) {
	tmp := t.TempDir()
	recap.SetHomeDir(t, func() (string, error) { return tmp, nil })

	recap.Save(&recap.Recap{Name: "dev", Description: "feature work"})
	recap.Save(&recap.Recap{Name: "build", Description: "CI pipeline"})

	msg := loadRecaps().(recapsMsg)
	if len(msg.recaps) != 2 {
		t.Fatalf("loadRecaps() returned %d recaps, want 2", len(msg.recaps))
	}
	if msg.recaps["dev"].Description != "feature work" {
		t.Errorf("dev recap = %q", msg.recaps["dev"].Description)
	}
}

func TestLoadRecapsEmpty(t *testing.T) {
	tmp := t.TempDir()
	recap.SetHomeDir(t, func() (string, error) { return tmp, nil })

	msg := loadRecaps().(recapsMsg)
	if len(msg.recaps) != 0 {
		t.Errorf("loadRecaps() returned %d recaps, want 0", len(msg.recaps))
	}
}

func TestRefreshSessionsListPanesError(t *testing.T) {
	tmux.SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|0|1\n"), nil })
	tmux.SetRunListPanes(t, func(name string) ([]byte, error) { return nil, fmt.Errorf("pane error") })
	process.SetRunPS(t, func() ([]byte, error) {
		return []byte("  PID  PPID COMM\n  100     1 bash\n"), nil
	})

	sm := refreshSessions().(sessionsMsg)
	if sm.err != nil {
		t.Fatalf("unexpected error: %v", sm.err)
	}
	if sm.sessions[0].HasClaude {
		t.Error("should not have claude when panes fail")
	}
}
