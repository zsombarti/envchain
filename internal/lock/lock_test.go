package lock_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"envchain/internal/lock"
)

func newManager(t *testing.T) *lock.Manager {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "locks")
	return lock.NewManager(dir)
}

func TestAcquireAndRelease(t *testing.T) {
	m := newManager(t)
	if err := m.Acquire("dev", "test"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !m.IsLocked("dev") {
		t.Fatal("expected profile to be locked")
	}
	if err := m.Release("dev"); err != nil {
		t.Fatalf("Release: %v", err)
	}
	if m.IsLocked("dev") {
		t.Fatal("expected profile to be unlocked after release")
	}
}

func TestAcquireDuplicate(t *testing.T) {
	m := newManager(t)
	_ = m.Acquire("dev", "first")
	err := m.Acquire("dev", "second")
	if err != lock.ErrLocked {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
}

func TestReleaseNotLocked(t *testing.T) {
	m := newManager(t)
	err := m.Release("dev")
	if err != lock.ErrNotLocked {
		t.Fatalf("expected ErrNotLocked, got %v", err)
	}
}

func TestStatusInfo(t *testing.T) {
	m := newManager(t)
	before := time.Now().UTC().Truncate(time.Second)
	if err := m.Acquire("staging", "alice"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	info, err := m.Status("staging")
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if info.Owner != "alice" {
		t.Errorf("owner = %q, want %q", info.Owner, "alice")
	}
	if info.AcquiredAt.Before(before) {
		t.Errorf("AcquiredAt %v is before test start %v", info.AcquiredAt, before)
	}
}

func TestStatusNotLocked(t *testing.T) {
	m := newManager(t)
	_, err := m.Status("ghost")
	if err != lock.ErrNotLocked {
		t.Fatalf("expected ErrNotLocked, got %v", err)
	}
}

func TestAcquireCreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "locks")
	m := lock.NewManager(dir)
	if err := m.Acquire("prod", "ci"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected dir to be created: %v", err)
	}
}

func TestIsLockedFalseWhenFree(t *testing.T) {
	m := newManager(t)
	if m.IsLocked("free") {
		t.Fatal("expected unlocked profile to report false")
	}
}
