package store_test

import (
	"errors"
	"testing"

	"envchain/internal/store"
)

func TestMemoryStoreSetAndGet(t *testing.T) {
	s := store.NewMemoryStore()
	if err := s.Set("API_KEY", "abc123"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	v, err := s.Get("API_KEY")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if v != "abc123" {
		t.Errorf("expected abc123, got %s", v)
	}
}

func TestMemoryStoreGetNotFound(t *testing.T) {
	s := store.NewMemoryStore()
	_, err := s.Get("MISSING")
	if !errors.Is(err, store.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	s := store.NewMemoryStore()
	_ = s.Set("TOKEN", "secret")
	if err := s.Delete("TOKEN"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err := s.Get("TOKEN")
	if !errors.Is(err, store.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound after delete, got %v", err)
	}
}

func TestMemoryStoreDeleteNotFound(t *testing.T) {
	s := store.NewMemoryStore()
	err := s.Delete("NONEXISTENT")
	if !errors.Is(err, store.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestMemoryStoreList(t *testing.T) {
	s := store.NewMemoryStore()
	_ = s.Set("A", "1")
	_ = s.Set("B", "2")
	keys, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestMemoryStoreListEmpty(t *testing.T) {
	s := store.NewMemoryStore()
	keys, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}
