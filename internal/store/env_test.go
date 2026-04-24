package store_test

import (
	"errors"
	"os"
	"testing"

	"envchain/internal/store"
)

func TestEnvStoreGet(t *testing.T) {
	const key = "ENVCHAIN_TEST_VAR"
	t.Setenv(key, "hello")
	s := store.NewEnvStore("")
	v, err := s.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if v != "hello" {
		t.Errorf("expected hello, got %s", v)
	}
}

func TestEnvStoreGetMissing(t *testing.T) {
	s := store.NewEnvStore("")
	_ = os.Unsetenv("ENVCHAIN_DEFINITELY_MISSING")
	_, err := s.Get("ENVCHAIN_DEFINITELY_MISSING")
	if !errors.Is(err, store.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestEnvStoreSetReadOnly(t *testing.T) {
	s := store.NewEnvStore("")
	if err := s.Set("X", "y"); err == nil {
		t.Error("expected error from Set on EnvStore")
	}
}

func TestEnvStoreDeleteReadOnly(t *testing.T) {
	s := store.NewEnvStore("")
	if err := s.Delete("X"); err == nil {
		t.Error("expected error from Delete on EnvStore")
	}
}

func TestEnvStoreListWithPrefix(t *testing.T) {
	t.Setenv("ENVCHAIN_ALPHA", "1")
	t.Setenv("ENVCHAIN_BETA", "2")
	s := store.NewEnvStore("ENVCHAIN_")
	keys, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	matched := 0
	for _, k := range keys {
		if k == "ENVCHAIN_ALPHA" || k == "ENVCHAIN_BETA" {
			matched++
		}
	}
	if matched < 2 {
		t.Errorf("expected at least 2 matching keys, got %d", matched)
	}
}
