package diff_test

import (
	"testing"

	"envchain/internal/diff"
)

func TestComputeAdded(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	r := diff.Compute(old, next)
	if len(r.Added()) != 1 || r.Added()[0].Key != "B" {
		t.Fatalf("expected B to be added, got %+v", r.Added())
	}
}

func TestComputeRemoved(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	r := diff.Compute(old, next)
	if len(r.Removed()) != 1 || r.Removed()[0].Key != "B" {
		t.Fatalf("expected B to be removed, got %+v", r.Removed())
	}
}

func TestComputeModified(t *testing.T) {
	old := map[string]string{"A": "old"}
	next := map[string]string{"A": "new"}
	r := diff.Compute(old, next)
	m := r.Modified()
	if len(m) != 1 || m[0].OldValue != "old" || m[0].NewValue != "new" {
		t.Fatalf("unexpected modified: %+v", m)
	}
}

func TestComputeUnchanged(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1"}
	r := diff.Compute(old, next)
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestHasChanges(t *testing.T) {
	old := map[string]string{}
	next := map[string]string{"X": "y"}
	if !diff.Compute(old, next).HasChanges() {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestComputeEmpty(t *testing.T) {
	r := diff.Compute(map[string]string{}, map[string]string{})
	if len(r.Changes) != 0 {
		t.Fatalf("expected empty diff, got %+v", r.Changes)
	}
}

func TestComputeSortedKeys(t *testing.T) {
	old := map[string]string{"Z": "1", "A": "1", "M": "1"}
	next := map[string]string{"Z": "2", "A": "1", "M": "2"}
	r := diff.Compute(old, next)
	keys := make([]string, len(r.Changes))
	for i, c := range r.Changes {
		keys[i] = c.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("keys not sorted: %v", keys)
	}
}
