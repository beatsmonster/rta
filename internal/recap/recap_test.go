package recap

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setTestHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	SetHomeDir(t, func() (string, error) { return tmp, nil })
	return tmp
}

func TestDir(t *testing.T) {
	tmp := setTestHome(t)
	dir := Dir()
	want := filepath.Join(tmp, ".rta", "recaps")
	if dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Dir() did not create directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("Dir() path is not a directory")
	}
}

func TestDirHomeDirError(t *testing.T) {
	SetHomeDir(t, func() (string, error) { return "", os.ErrPermission })
	dir := Dir()
	if !filepath.IsAbs(dir) && dir != filepath.Join(".", ".rta", "recaps") {
		t.Logf("Dir() with error = %q (using fallback)", dir)
	}
}

func TestSaveAndLoad(t *testing.T) {
	setTestHome(t)
	now := time.Now().Truncate(time.Second)
	r := &Recap{
		Name:        "dev",
		Description: "working on feature X",
		ProjectPath: "/home/user/project",
		SessionID:   "abc123",
		UpdatedAt:   now,
	}
	if err := Save(r); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	loaded, err := Load("dev")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("Load() returned nil")
	}
	if loaded.Name != "dev" {
		t.Errorf("Name = %q, want %q", loaded.Name, "dev")
	}
	if loaded.Description != "working on feature X" {
		t.Errorf("Description = %q, want %q", loaded.Description, "working on feature X")
	}
	if loaded.ProjectPath != "/home/user/project" {
		t.Errorf("ProjectPath = %q", loaded.ProjectPath)
	}
	if loaded.SessionID != "abc123" {
		t.Errorf("SessionID = %q", loaded.SessionID)
	}
	if !loaded.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", loaded.UpdatedAt, now)
	}
}

func TestLoadMissingFile(t *testing.T) {
	setTestHome(t)
	r, err := Load("nonexistent")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if r != nil {
		t.Errorf("Load() = %+v, want nil", r)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmp := setTestHome(t)
	dir := filepath.Join(tmp, ".rta", "recaps")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "bad.json"), []byte("{invalid"), 0o644)

	_, err := Load("bad")
	if err == nil {
		t.Error("Load() should return error for invalid JSON")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	tmp := setTestHome(t)
	r := &Recap{Name: "test", Description: "desc"}
	if err := Save(r); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	path := filepath.Join(tmp, ".rta", "recaps", "test.json")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("Save() did not create file: %v", err)
	}
}

func TestList(t *testing.T) {
	setTestHome(t)
	r1 := &Recap{Name: "alpha", Description: "first"}
	r2 := &Recap{Name: "beta", Description: "second"}
	Save(r1)
	Save(r2)

	recaps, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(recaps) != 2 {
		t.Fatalf("List() returned %d recaps, want 2", len(recaps))
	}
}

func TestListEmpty(t *testing.T) {
	setTestHome(t)
	recaps, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(recaps) != 0 {
		t.Errorf("List() returned %d recaps, want 0", len(recaps))
	}
}

func TestListSkipsNonJSON(t *testing.T) {
	tmp := setTestHome(t)
	dir := filepath.Join(tmp, ".rta", "recaps")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not json"), 0o644)
	os.MkdirAll(filepath.Join(dir, "subdir"), 0o755)
	Save(&Recap{Name: "real", Description: "valid"})

	recaps, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(recaps) != 1 {
		t.Errorf("List() returned %d recaps, want 1", len(recaps))
	}
}

func TestLoadUnreadableFile(t *testing.T) {
	tmp := setTestHome(t)
	dir := filepath.Join(tmp, ".rta", "recaps")
	os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, "noperm.json")
	os.WriteFile(path, []byte(`{"name":"noperm"}`), 0o644)
	os.Chmod(path, 0o000)
	t.Cleanup(func() { os.Chmod(path, 0o644) })

	_, err := Load("noperm")
	if err == nil {
		t.Error("Load() should return error for unreadable file")
	}
}

func TestListSkipsInvalidJSON(t *testing.T) {
	tmp := setTestHome(t)
	dir := filepath.Join(tmp, ".rta", "recaps")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "broken.json"), []byte("{bad"), 0o644)
	Save(&Recap{Name: "good", Description: "valid"})

	recaps, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(recaps) != 1 {
		t.Errorf("List() returned %d recaps, want 1", len(recaps))
	}
}
