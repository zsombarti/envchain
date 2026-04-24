package lock_test

import (
	"errors"
	"testing"

	"envchain/internal/lock"
)

func TestWithLockAcquiresAndReleases(t *testing.T) {
	m := newManager(t)
	guard, err := lock.WithLock(m, "dev", "tester")
	if err != nil {
		t.Fatalf("WithLock: %v", err)
	}
	if !m.IsLocked("dev") {
		t.Fatal("expected lock to be held")
	}
	if err := guard.Done(); err != nil {
		t.Fatalf("Done: %v", err)
	}
	if m.IsLocked("dev") {
		t.Fatal("expected lock to be released after Done")
	}
}

func TestWithLockReturnsErrLocked(t *testing.T) {
	m := newManager(t)
	guard, err := lock.WithLock(m, "dev", "first")
	if err != nil {
		t.Fatalf("first WithLock: %v", err)
	}
	defer guard.Done()

	_, err = lock.WithLock(m, "dev", "second")
	if !errors.Is(err, lock.ErrLocked) {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
}

func TestGuardDoneIdempotent(t *testing.T) {
	m := newManager(t)
	guard, err := lock.WithLock(m, "dev", "tester")
	if err != nil {
		t.Fatalf("WithLock: %v", err)
	}
	if err := guard.Done(); err != nil {
		t.Fatalf("first Done: %v", err)
	}
	// Second Done should not return an error (ErrNotLocked is swallowed).
	if err := guard.Done(); err != nil {
		t.Fatalf("second Done: %v", err)
	}
}

func TestWithLockDifferentProfiles(t *testing.T) {
	m := newManager(t)
	g1, err := lock.WithLock(m, "dev", "a")
	if err != nil {
		t.Fatalf("lock dev: %v", err)
	}
	g2, err := lock.WithLock(m, "staging", "b")
	if err != nil {
		t.Fatalf("lock staging: %v", err)
	}
	defer g1.Done()
	defer g2.Done()
	if !m.IsLocked("dev") || !m.IsLocked("staging") {
		t.Fatal("both profiles should be locked independently")
	}
}
