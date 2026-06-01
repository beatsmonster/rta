package tmux

import "testing"

func SetRunListSessions(t *testing.T, fn func() ([]byte, error)) {
	orig := runListSessions
	t.Cleanup(func() { runListSessions = orig })
	runListSessions = fn
}

func SetRunListPanes(t *testing.T, fn func(string) ([]byte, error)) {
	orig := runListPanes
	t.Cleanup(func() { runListPanes = orig })
	runListPanes = fn
}

func SetLookPath(t *testing.T, fn func(string) (string, error)) {
	orig := lookPath
	t.Cleanup(func() { lookPath = orig })
	lookPath = fn
}

func SetExecSyscall(t *testing.T, fn func(string, []string, []string) error) {
	orig := execSyscall
	t.Cleanup(func() { execSyscall = orig })
	execSyscall = fn
}
