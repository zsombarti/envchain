package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"envchain/internal/profile"
)

func newManager(t *testing.T) *profile.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "profiles")
	return profile.NewManager(dir)
}

func TestSaveAndLoad(t *testing.T) {
	m := newManager(t)
	p := profile.Profile{Name: "dev", Layers: []string{"base", "local"}}
	if err := m.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != p.Name {
		t.Errorf("name: got %q, want %q", got.Name, p.Name)
	}
	if len(got.Layers) != len(p.Layers) {
		t.Errorf("layers len: got %d, want %d", len(got.Layers), len(p.Layers))
	}
}

func TestLoadNotFound(t *testing.T) {
	m := newManager(t)
	_, err := m.Load("missing")
	if err != profile.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newManager(t)
	p := profile.Profile{Name: "staging", Layers: []string{"base"}}
	_ = m.Save(p)
	if err := m.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Load("staging")
	if err != profile.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newManager(t)
	if err := m.Delete("ghost"); err != profile.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := newManager(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		_ = m.Save(profile.Profile{Name: name, Layers: []string{"base"}})
	}
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(names))
	}
}

func TestListEmpty(t *testing.T) {
	m := newManager(t)
	names, err := m.List()
	if err != nil {
		t.Fatalf("List on empty dir: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(names))
	}
}

func TestSaveCreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "profiles")
	m := profile.NewManager(dir)
	p := profile.Profile{Name: "ci", Layers: []string{}}
	if err := m.Save(p); err != nil {
		t.Fatalf("Save in nested dir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}
