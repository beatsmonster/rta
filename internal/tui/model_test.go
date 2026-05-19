package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"rta/internal/recap"
)

type testError struct{}

func (e *testError) Error() string { return "test error" }

func TestNew(t *testing.T) {
	m := New()
	if m.cursor != 0 || len(m.sessions) != 0 {
		t.Errorf("New() = %+v, want zero value", m)
	}
}

func TestInit(t *testing.T) {
	m := New()
	if cmd := m.Init(); cmd == nil {
		t.Fatal("Init() returned nil cmd")
	}
}

func TestUpdateKeyUp(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "a"}, {Name: "b"}}, cursor: 1}
	// Code: -1 means no physical key code; bubbletea uses the Text field for matching.
	um, _ := m.Update(tea.KeyPressMsg{Code: -1, Text: "k"})
	if um.(model).cursor != 0 {
		t.Errorf("cursor = %d, want 0", um.(model).cursor)
	}
}

func TestUpdateKeyUpAtZero(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "a"}}, cursor: 0}
	um, _ := m.Update(tea.KeyPressMsg{Code: -1, Text: "up"})
	if um.(model).cursor != 0 {
		t.Errorf("cursor = %d, want 0", um.(model).cursor)
	}
}

func TestUpdateKeyDown(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "a"}, {Name: "b"}}, cursor: 0}
	um, _ := m.Update(tea.KeyPressMsg{Code: -1, Text: "j"})
	if um.(model).cursor != 1 {
		t.Errorf("cursor = %d, want 1", um.(model).cursor)
	}
}

func TestUpdateKeyDownAtEnd(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "a"}, {Name: "b"}}, cursor: 1}
	um, _ := m.Update(tea.KeyPressMsg{Code: -1, Text: "down"})
	if um.(model).cursor != 1 {
		t.Errorf("cursor = %d, want 1", um.(model).cursor)
	}
}

func TestUpdateKeyEnterWithSessions(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "test"}}, cursor: 0}
	_, cmd := m.Update(tea.KeyPressMsg{Code: -1, Text: "enter"})
	if cmd == nil {
		t.Error("enter should return a cmd when sessions exist")
	}
}

func TestUpdateKeyEnterNoSessions(t *testing.T) {
	m := model{}
	_, cmd := m.Update(tea.KeyPressMsg{Code: -1, Text: "enter"})
	if cmd != nil {
		t.Error("enter should return nil cmd when no sessions")
	}
}

func TestUpdateKeyRefresh(t *testing.T) {
	_, cmd := model{}.Update(tea.KeyPressMsg{Code: -1, Text: "r"})
	if cmd == nil {
		t.Error("r should return refreshSessions cmd")
	}
}

func TestUpdateQuitKeys(t *testing.T) {
	for _, key := range []string{"q", "ctrl+c", "esc"} {
		t.Run(key, func(t *testing.T) {
			_, cmd := model{}.Update(tea.KeyPressMsg{Code: -1, Text: key})
			if cmd == nil {
				t.Errorf("%s should return quit cmd", key)
			}
		})
	}
}

func TestUpdateWindowSize(t *testing.T) {
	um, _ := model{}.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m := um.(model)
	if m.width != 120 || m.height != 40 {
		t.Errorf("size = %dx%d, want 120x40", m.width, m.height)
	}
}

func TestUpdateSessionsMsg(t *testing.T) {
	sessions := []SessionInfo{{Name: "dev"}, {Name: "build"}}
	um, _ := model{}.Update(sessionsMsg{sessions: sessions})
	m := um.(model)
	if len(m.sessions) != 2 || m.err != nil {
		t.Errorf("sessions=%d err=%v", len(m.sessions), m.err)
	}
}

func TestUpdateSessionsMsgError(t *testing.T) {
	um, _ := model{sessions: []SessionInfo{{Name: "old"}}}.Update(sessionsMsg{err: &testError{}})
	m := um.(model)
	if m.err == nil || m.sessions != nil {
		t.Error("error msg should set err and clear sessions")
	}
}

func TestUpdateSessionsMsgCursorClamp(t *testing.T) {
	um, _ := model{cursor: 5}.Update(sessionsMsg{sessions: []SessionInfo{{Name: "only"}}})
	if um.(model).cursor != 0 {
		t.Errorf("cursor = %d, want 0", um.(model).cursor)
	}
}

func TestUpdateSessionsMsgEmptyCursorClamp(t *testing.T) {
	um, _ := model{cursor: 5}.Update(sessionsMsg{sessions: nil})
	if um.(model).cursor != 0 {
		t.Errorf("cursor = %d, want 0", um.(model).cursor)
	}
}

func TestUpdateTmuxFinished(t *testing.T) {
	_, cmd := model{}.Update(tmuxFinishedMsg{})
	if cmd == nil {
		t.Error("tmuxFinishedMsg should return refresh cmd")
	}
}

func TestUpdateUnknownMsg(t *testing.T) {
	_, cmd := model{}.Update("unknown")
	if cmd != nil {
		t.Error("unknown msg should return nil cmd")
	}
}

func TestViewError(t *testing.T) {
	v := model{err: &testError{}}.View()
	if !strings.Contains(v.Content, "Error") || !strings.Contains(v.Content, "esc/q") {
		t.Errorf("error view missing content: %s", v.Content)
	}
}

func TestViewNoSessions(t *testing.T) {
	v := model{}.View()
	if !strings.Contains(v.Content, "No tmux sessions") || !strings.Contains(v.Content, "esc/q") {
		t.Errorf("empty view missing content: %s", v.Content)
	}
}

func TestViewWithClaudeSessions(t *testing.T) {
	m := model{sessions: []SessionInfo{
		{Name: "dev", HasClaude: true, WorkingDir: "/tmp/project"},
		{Name: "build", HasClaude: false, WorkingDir: "/tmp/build"},
	}}
	v := m.View()
	if !strings.Contains(v.Content, "Claude Code Sessions") || !strings.Contains(v.Content, "Other tmux Sessions") {
		t.Errorf("view missing sections: %s", v.Content)
	}
}

func TestViewOnlyOtherSessions(t *testing.T) {
	v := model{sessions: []SessionInfo{{Name: "build", HasClaude: false, WorkingDir: "/tmp"}}}.View()
	if strings.Contains(v.Content, "Claude Code Sessions") {
		t.Error("should not show claude section")
	}
}

func TestRenderSessionCursor(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "test", WorkingDir: "/tmp"}}, cursor: 0}
	if s := m.renderSession(0); !strings.Contains(s, ">") {
		t.Errorf("cursor session should have >, got: %s", s)
	}
}

func TestRenderSessionNotCursor(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "a", WorkingDir: "/tmp"}, {Name: "b", WorkingDir: "/tmp"}}, cursor: 0}
	if s := m.renderSession(1); strings.Contains(s, ">") {
		t.Errorf("non-cursor session should not have >, got: %s", s)
	}
}

func TestRenderSessionAttached(t *testing.T) {
	m := model{sessions: []SessionInfo{{Name: "test", WorkingDir: "/tmp", Attached: true}}, cursor: 0}
	if s := m.renderSession(0); !strings.Contains(s, "(attached)") {
		t.Errorf("attached session should show (attached), got: %s", s)
	}
}

func TestUpdateRecapsMsg(t *testing.T) {
	recaps := map[string]*recap.Recap{
		"dev": {Name: "dev", Description: "working on feature X"},
	}
	um, _ := model{}.Update(recapsMsg{recaps: recaps})
	m := um.(model)
	if m.recaps == nil || m.recaps["dev"] == nil {
		t.Error("recapsMsg should set recaps map")
	}
}

func TestViewWithRecap(t *testing.T) {
	m := model{
		sessions: []SessionInfo{{Name: "dev", HasClaude: true, WorkingDir: "/tmp"}},
		recaps: map[string]*recap.Recap{
			"dev": {Name: "dev", Description: "working on feature X"},
		},
		cursor: 0,
	}
	v := m.View()
	if !strings.Contains(v.Content, "working on feature X") {
		t.Errorf("view should show recap description, got: %s", v.Content)
	}
}

func TestViewWithoutRecap(t *testing.T) {
	m := model{
		sessions: []SessionInfo{{Name: "dev", HasClaude: true, WorkingDir: "/tmp"}},
		recaps:   map[string]*recap.Recap{},
		cursor:   0,
	}
	v := m.View()
	if strings.Contains(v.Content, "no recap") {
		t.Error("view should not show 'no recap' placeholder")
	}
}

func TestViewRecapNilMap(t *testing.T) {
	m := model{
		sessions: []SessionInfo{{Name: "dev", HasClaude: true, WorkingDir: "/tmp"}},
		recaps:   nil,
		cursor:   0,
	}
	v := m.View()
	if !strings.Contains(v.Content, "dev") {
		t.Errorf("view should still render sessions with nil recaps map, got: %s", v.Content)
	}
}

func TestViewRecapChangesWithCursor(t *testing.T) {
	m := model{
		sessions: []SessionInfo{
			{Name: "dev", HasClaude: true, WorkingDir: "/tmp"},
			{Name: "build", HasClaude: false, WorkingDir: "/tmp"},
		},
		recaps: map[string]*recap.Recap{
			"dev":   {Name: "dev", Description: "feature work"},
			"build": {Name: "build", Description: "CI pipeline"},
		},
		cursor: 0,
	}
	v0 := m.View()
	if !strings.Contains(v0.Content, "feature work") {
		t.Error("cursor=0 should show dev recap")
	}

	m.cursor = 1
	v1 := m.View()
	if !strings.Contains(v1.Content, "CI pipeline") {
		t.Error("cursor=1 should show build recap")
	}
	if strings.Contains(v1.Content, "feature work") {
		t.Error("cursor=1 should not show dev recap")
	}
}

func TestViewRecapEmptyDescription(t *testing.T) {
	m := model{
		sessions: []SessionInfo{{Name: "dev", HasClaude: true, WorkingDir: "/tmp"}},
		recaps: map[string]*recap.Recap{
			"dev": {Name: "dev", Description: ""},
		},
		cursor: 0,
	}
	v := m.View()
	lines := strings.Split(v.Content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.Contains(trimmed, "dev") || strings.Contains(trimmed, "Claude") ||
			strings.Contains(trimmed, "attach") || strings.Contains(trimmed, "Other") {
			continue
		}
	}
	_ = v
}
