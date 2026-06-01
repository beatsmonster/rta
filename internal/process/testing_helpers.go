package process

import "testing"

func SetRunPS(t *testing.T, fn func() ([]byte, error)) {
	orig := runPS
	t.Cleanup(func() { runPS = orig })
	runPS = fn
}
