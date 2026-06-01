package recap

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Recap struct {
	Name        string    `json:"name"`
	Description string    `json:"recap"`
	ProjectPath string    `json:"project_path"`
	SessionID   string    `json:"session_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var getHomeDir = os.UserHomeDir

func Dir() string {
	home, err := getHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".rta", "recaps")
	os.MkdirAll(dir, 0o755)
	return dir
}

func Load(name string) (*Recap, error) {
	path := filepath.Join(Dir(), name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var r Recap
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func Save(r *Recap) error {
	dir := Dir()
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, r.Name+".json"), data, 0o644)
}

func List() ([]*Recap, error) {
	entries, err := os.ReadDir(Dir())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var recaps []*Recap
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		r, err := Load(name)
		if err != nil {
			continue
		}
		if r != nil {
			recaps = append(recaps, r)
		}
	}
	return recaps, nil
}
