package tmux

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestParseSessions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Session
		wantLen int
	}{
		{
			name:    "two sessions",
			input:   "webapp|1|3\nbuild|0|1\n",
			want:    []Session{{Name: "webapp", Attached: true, Windows: 3}, {Name: "build", Attached: false, Windows: 1}},
			wantLen: 2,
		},
		{
			name:    "single session",
			input:   "dev|0|2\n",
			want:    []Session{{Name: "dev", Attached: false, Windows: 2}},
			wantLen: 1,
		},
		{
			name:    "empty output",
			input:   "",
			want:    nil,
			wantLen: 0,
		},
		{
			name:    "whitespace only",
			input:   "  \n  \n",
			want:    nil,
			wantLen: 0,
		},
		{
			name:    "malformed line skipped",
			input:   "webapp|1|3\nbadline\nbuild|0|1\n",
			want:    []Session{{Name: "webapp", Attached: true, Windows: 3}, {Name: "build", Attached: false, Windows: 1}},
			wantLen: 2,
		},
		{
			name:    "session name with spaces",
			input:   "my session|1|2\n",
			want:    []Session{{Name: "my session", Attached: true, Windows: 2}},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSessions(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("got %d sessions, want %d", len(got), tt.wantLen)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("session[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParsePanes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Pane
		wantLen int
	}{
		{
			name:    "two panes",
			input:   "12345|/Users/alice/project\n12400|/Users/alice/other\n",
			want:    []Pane{{PID: 12345, CurrentPath: "/Users/alice/project"}, {PID: 12400, CurrentPath: "/Users/alice/other"}},
			wantLen: 2,
		},
		{
			name:    "path with pipe uses SplitN",
			input:   "99999|/some/path|with|pipes\n",
			want:    []Pane{{PID: 99999, CurrentPath: "/some/path|with|pipes"}},
			wantLen: 1,
		},
		{
			name:    "empty output",
			input:   "",
			want:    nil,
			wantLen: 0,
		},
		{
			name:    "non-numeric pid skipped",
			input:   "abc|/some/path\n12345|/good/path\n",
			want:    []Pane{{PID: 12345, CurrentPath: "/good/path"}},
			wantLen: 1,
		},
		{
			name:    "missing fields skipped",
			input:   "12345\n",
			want:    nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePanes(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("got %d panes, want %d", len(got), tt.wantLen)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("pane[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFindSession(t *testing.T) {
	sessions := []Session{
		{Name: "webapp", Attached: true, Windows: 3},
		{Name: "web-api", Attached: false, Windows: 1},
		{Name: "build", Attached: false, Windows: 2},
	}

	tests := []struct {
		name      string
		substring string
		sessions  []Session
		wantNames []string
	}{
		{
			name:      "exact match",
			substring: "webapp",
			sessions:  sessions,
			wantNames: []string{"webapp"},
		},
		{
			name:      "substring matches multiple",
			substring: "web",
			sessions:  sessions,
			wantNames: []string{"webapp", "web-api"},
		},
		{
			name:      "no match",
			substring: "nonexistent",
			sessions:  sessions,
			wantNames: nil,
		},
		{
			name:      "empty sessions",
			substring: "web",
			sessions:  nil,
			wantNames: nil,
		},
		{
			name:      "empty substring matches all",
			substring: "",
			sessions:  sessions,
			wantNames: []string{"webapp", "web-api", "build"},
		},
		{
			name:      "single character match",
			substring: "b",
			sessions:  sessions,
			wantNames: []string{"webapp", "web-api", "build"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSession(tt.substring, tt.sessions)
			if len(got) != len(tt.wantNames) {
				t.Fatalf("got %d matches, want %d", len(got), len(tt.wantNames))
			}
			for i, s := range got {
				if s.Name != tt.wantNames[i] {
					t.Errorf("match[%d].Name = %q, want %q", i, s.Name, tt.wantNames[i])
				}
			}
		})
	}
}

func TestListSessionsSuccess(t *testing.T) {
	SetRunListSessions(t, func() ([]byte, error) { return []byte("dev|1|2\nbuild|0|1\n"), nil })

	sessions, err := ListSessions()
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	if len(sessions) != 2 || sessions[0].Name != "dev" || !sessions[0].Attached {
		t.Errorf("unexpected sessions: %+v", sessions)
	}
}

func TestListSessionsError(t *testing.T) {
	SetRunListSessions(t, func() ([]byte, error) {
		return nil, &exec.ExitError{Stderr: []byte("no server running")}
	})

	if _, err := ListSessions(); err != ErrNoServer {
		t.Errorf("got %v, want ErrNoServer", err)
	}
}

func TestListPanesSuccess(t *testing.T) {
	SetRunListPanes(t, func(name string) ([]byte, error) { return []byte("12345|/home/user/project\n"), nil })

	panes, err := ListPanes("dev")
	if err != nil {
		t.Fatalf("ListPanes: %v", err)
	}
	if len(panes) != 1 || panes[0].PID != 12345 {
		t.Errorf("unexpected panes: %+v", panes)
	}
}

func TestListPanesError(t *testing.T) {
	SetRunListPanes(t, func(name string) ([]byte, error) { return nil, fmt.Errorf("tmux error") })

	if _, err := ListPanes("dev"); err == nil {
		t.Error("expected error")
	}
}

func TestAttachSessionSuccess(t *testing.T) {
	SetLookPath(t, func(file string) (string, error) { return "/usr/bin/tmux", nil })
	SetExecSyscall(t, func(binary string, args []string, env []string) error {
		if binary != "/usr/bin/tmux" || args[3] != "test-session" {
			t.Errorf("unexpected args: binary=%q args=%v", binary, args)
		}
		return nil
	})

	if err := AttachSession("test-session"); err != nil {
		t.Fatalf("AttachSession: %v", err)
	}
}

func TestAttachSessionLookPathError(t *testing.T) {
	SetLookPath(t, func(file string) (string, error) { return "", fmt.Errorf("not found") })

	err := AttachSession("test")
	if err == nil || !strings.Contains(err.Error(), "tmux not found") {
		t.Errorf("expected 'tmux not found', got: %v", err)
	}
}

func TestClassifyErrorNonExitError(t *testing.T) {
	err := classifyError(fmt.Errorf("generic error"))
	if err == ErrNoServer || err == ErrNoSessions {
		t.Errorf("expected generic wrapped error, got %v", err)
	}
}

func TestClassifyError_NoServer(t *testing.T) {
	err := &exec.ExitError{Stderr: []byte("error connecting to /tmp/tmux-501/default (no server running on /tmp/tmux-501/default)")}
	got := classifyError(err)
	if got != ErrNoServer {
		t.Errorf("got %v, want ErrNoServer", got)
	}
}

func TestClassifyError_NoSessions(t *testing.T) {
	err := &exec.ExitError{Stderr: []byte("no sessions")}
	got := classifyError(err)
	if got != ErrNoSessions {
		t.Errorf("got %v, want ErrNoSessions", got)
	}
}

func TestClassifyError_Other(t *testing.T) {
	err := &exec.ExitError{Stderr: []byte("some other error")}
	got := classifyError(err)
	if got == ErrNoServer || got == ErrNoSessions {
		t.Errorf("expected generic error, got %v", got)
	}
}
