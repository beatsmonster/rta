package profile

import "testing"

func SetRunPgrep(t *testing.T, fn func() ([]byte, error)) {
	orig := runPgrep
	t.Cleanup(func() { runPgrep = orig })
	runPgrep = fn
}

func SetGetHomeDir(t *testing.T, fn func() (string, error)) {
	orig := getHomeDir
	t.Cleanup(func() { getHomeDir = orig })
	getHomeDir = fn
}
