package recap

import "testing"

func SetHomeDir(t *testing.T, fn func() (string, error)) {
	orig := getHomeDir
	t.Cleanup(func() { getHomeDir = orig })
	getHomeDir = fn
}
