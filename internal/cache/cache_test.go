package cache

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSetAndGet(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("dev", map[string]string{"FOO": "bar"})
	got, ok := c.Get("dev")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["FOO"] != "bar" {
		t.Fatalf("expected bar, got %s", got["FOO"])
	}
}

func TestGetMiss(t *testing.T) {
	c := New(5 * time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestExpiredEntry(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c := New(1 * time.Minute)
	c.now = fixedNow(base)
	c.Set("dev", map[string]string{"A": "1"})

	// Advance past TTL
	c.now = fixedNow(base.Add(2 * time.Minute))
	_, ok := c.Get("dev")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestGetReturnsCopy(t *testing.T) {
	c := New(5 * time.Minute)
	orig := map[string]string{"X": "1"}
	c.Set("p", orig)
	got, _ := c.Get("p")
	got["X"] = "mutated"
	again, _ := c.Get("p")
	if again["X"] != "1" {
		t.Fatal("cache returned mutable reference")
	}
}

func TestDelete(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("dev", map[string]string{"K": "v"})
	c.Delete("dev")
	_, ok := c.Get("dev")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestDeleteNotFound(t *testing.T) {
	c := New(5 * time.Minute)
	// Should not panic
	c.Delete("nonexistent")
}

func TestFlush(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("a", map[string]string{"A": "1"})
	c.Set("b", map[string]string{"B": "2"})
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", c.Len())
	}
}

func TestLen(t *testing.T) {
	c := New(5 * time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("x", map[string]string{})
	c.Set("y", map[string]string{})
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
