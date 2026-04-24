package snapshot_test

import (
	"testing"

	"envchain/internal/snapshot"
)

func newManager(t *testing.T) *snapshot.Manager {
	t.Helper()
	dir := t.TempDir()
	m, err := snapshot.NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestManagerSaveAndLoad(t *testing.T) {
	m := newManager(t)
	s := snapshot.New("dev", map[string]string{"DB_URL": "postgres://localhost/dev"})

	if err := m.Save(s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := m.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Vars["DB_URL"] != "postgres://localhost/dev" {
		t.Errorf("unexpected DB_URL: %q", loaded.Vars["DB_URL"])
	}
}

func TestManagerLoadNotFound(t *testing.T) {
	m := newManager(t)
	_, err := m.Load("missing")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestManagerDelete(t *testing.T) {
	m := newManager(t)
	s := snapshot.New("to-delete", map[string]string{"X": "1"})
	if err := m.Save(s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Delete("to-delete"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Load("to-delete")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestManagerDeleteNotFound(t *testing.T) {
	m := newManager(t)
	if err := m.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent snapshot")
	}
}

func TestManagerList(t *testing.T) {
	m := newManager(t)
	for _, label := range []string{"alpha", "beta", "gamma"} {
		if err := m.Save(snapshot.New(label, map[string]string{})); err != nil {
			t.Fatalf("Save %q: %v", label, err)
		}
	}
	labels, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 3 {
		t.Errorf("expected 3 labels, got %d: %v", len(labels), labels)
	}
}

func TestManagerListEmpty(t *testing.T) {
	m := newManager(t)
	labels, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty list, got %v", labels)
	}
}
