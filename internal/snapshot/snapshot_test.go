package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"envchain/internal/snapshot"
)

func TestNewSnapshot(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New("test", vars)

	if s.Label != "test" {
		t.Errorf("expected label 'test', got %q", s.Label)
	}
	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	// Ensure it's a copy
	vars["EXTRA"] = "val"
	if _, ok := s.Vars["EXTRA"]; ok {
		t.Error("snapshot should not reflect mutations to source map")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New("save-load", map[string]string{"KEY": "value"})
	orig.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Label != orig.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, orig.Label)
	}
	if loaded.Vars["KEY"] != "value" {
		t.Errorf("var mismatch: got %q", loaded.Vars["KEY"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error loading missing file")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	s := snapshot.New("perm-test", map[string]string{})
	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}

func TestDiff(t *testing.T) {
	a := snapshot.New("a", map[string]string{"SAME": "v", "CHANGED": "old", "REMOVED": "gone"})
	b := snapshot.New("b", map[string]string{"SAME": "v", "CHANGED": "new", "ADDED": "here"})

	added, removed, changed := snapshot.Diff(a, b)

	if len(added) != 1 || added[0] != "ADDED" {
		t.Errorf("added: expected [ADDED], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "REMOVED" {
		t.Errorf("removed: expected [REMOVED], got %v", removed)
	}
	if len(changed) != 1 || changed[0] != "CHANGED" {
		t.Errorf("changed: expected [CHANGED], got %v", changed)
	}
}
