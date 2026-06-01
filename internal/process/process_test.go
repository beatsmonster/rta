package process

import (
	"fmt"
	"testing"
)

const mockPSOutput = `  PID  PPID COMM
    1     0 launchd
  256     1 sshd
  257   256 bash
  300   257 claude
  400     1 WindowServer
  500   400 Finder
`

func TestParseTree(t *testing.T) {
	children, names, err := parseTree([]byte(mockPSOutput))
	if err != nil {
		t.Fatalf("parseTree: %v", err)
	}

	if names[1] != "launchd" {
		t.Errorf("names[1] = %q, want launchd", names[1])
	}
	if names[300] != "claude" {
		t.Errorf("names[300] = %q, want claude", names[300])
	}

	if got := children[256]; len(got) != 1 || got[0] != 257 {
		t.Errorf("children[256] = %v, want [257]", got)
	}
	if got := children[1]; len(got) != 2 {
		t.Errorf("children[1] has %d entries, want 2", len(got))
	}
}

func TestParseTreeEmpty(t *testing.T) {
	children, names, err := parseTree([]byte("  PID  PPID COMM\n"))
	if err != nil {
		t.Fatalf("parseTree: %v", err)
	}
	if len(children) != 0 {
		t.Errorf("children = %v, want empty", children)
	}
	if len(names) != 0 {
		t.Errorf("names = %v, want empty", names)
	}
}

func TestParseTreeMalformed(t *testing.T) {
	input := "  PID  PPID COMM\n  abc  def  foo\n  100  200  bar\n"
	children, names, err := parseTree([]byte(input))
	if err != nil {
		t.Fatalf("parseTree: %v", err)
	}
	if len(names) != 1 || names[100] != "bar" {
		t.Errorf("expected only pid 100, got names=%v", names)
	}
	if got := children[200]; len(got) != 1 || got[0] != 100 {
		t.Errorf("children[200] = %v, want [100]", got)
	}
}

func TestHasDescendant(t *testing.T) {
	children, names, _ := parseTree([]byte(mockPSOutput))

	if !HasDescendant(256, "claude", children, names) {
		t.Error("expected claude descendant of 256")
	}
	if !HasDescendant(257, "claude", children, names) {
		t.Error("expected claude descendant of 257")
	}
	if HasDescendant(400, "claude", children, names) {
		t.Error("did not expect claude descendant of 400")
	}
}

func TestHasDescendantRootIsTarget(t *testing.T) {
	children, names, _ := parseTree([]byte(mockPSOutput))

	if !HasDescendant(300, "claude", children, names) {
		t.Error("expected root pid 300 to match claude")
	}
}

func TestHasDescendantMissingPID(t *testing.T) {
	children, names, _ := parseTree([]byte(mockPSOutput))

	if HasDescendant(99999, "claude", children, names) {
		t.Error("expected false for nonexistent PID")
	}
}

func TestBuildTreeWithMock(t *testing.T) {
	SetRunPS(t, func() ([]byte, error) {
		return []byte(mockPSOutput), nil
	})

	children, names, err := BuildTree()
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}
	if names[300] != "claude" {
		t.Errorf("names[300] = %q, want claude", names[300])
	}
	if len(children[1]) != 2 {
		t.Errorf("children[1] has %d entries, want 2", len(children[1]))
	}
}

func TestBuildTreeError(t *testing.T) {
	SetRunPS(t, func() ([]byte, error) { return nil, fmt.Errorf("ps failed") })

	if _, _, err := BuildTree(); err == nil {
		t.Error("expected error from BuildTree")
	}
}

func TestParseTreeShortFields(t *testing.T) {
	children, names, err := parseTree([]byte("  PID  PPID COMM\n  100  200\n  300  400  bash\n"))
	if err != nil {
		t.Fatalf("parseTree: %v", err)
	}
	if len(names) != 1 || names[300] != "bash" {
		t.Errorf("expected only pid 300, got names=%v", names)
	}
	if len(children[400]) != 1 {
		t.Errorf("children[400] = %v, want [300]", children[400])
	}
}

func TestDetectClaudeError(t *testing.T) {
	SetRunPS(t, func() ([]byte, error) { return nil, fmt.Errorf("ps failed") })

	if _, err := DetectClaude(1); err == nil {
		t.Error("expected error from DetectClaude")
	}
}

func TestDetectClaudeWithMock(t *testing.T) {
	SetRunPS(t, func() ([]byte, error) {
		return []byte(mockPSOutput), nil
	})

	found, err := DetectClaude(256)
	if err != nil {
		t.Fatalf("DetectClaude: %v", err)
	}
	if !found {
		t.Error("expected DetectClaude(256) = true")
	}

	found, err = DetectClaude(400)
	if err != nil {
		t.Fatalf("DetectClaude: %v", err)
	}
	if found {
		t.Error("expected DetectClaude(400) = false")
	}
}
