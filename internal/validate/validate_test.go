package validate_test

import (
	"errors"
	"testing"

	"envchain/internal/validate"
)

func TestKeyValid(t *testing.T) {
	valid := []string{
		"FOO",
		"_BAR",
		"MY_VAR_123",
		"a",
		"_",
		"Z9",
	}
	for _, k := range valid {
		if err := validate.Key(k); err != nil {
			t.Errorf("Key(%q) unexpected error: %v", k, err)
		}
	}
}

func TestKeyEmpty(t *testing.T) {
	if err := validate.Key(""); !errors.Is(err, validate.ErrEmptyKey) {
		t.Fatalf("expected ErrEmptyKey, got %v", err)
	}
}

func TestKeyInvalidStart(t *testing.T) {
	invalid := []string{"1FOO", "9", "-BAR", "=X"}
	for _, k := range invalid {
		err := validate.Key(k)
		if !errors.Is(err, validate.ErrInvalidKey) {
			t.Errorf("Key(%q): expected ErrInvalidKey, got %v", k, err)
		}
	}
}

func TestKeyInvalidCharacter(t *testing.T) {
	invalid := []string{"FOO-BAR", "MY.VAR", "A B", "X@Y"}
	for _, k := range invalid {
		err := validate.Key(k)
		if !errors.Is(err, validate.ErrInvalidKey) {
			t.Errorf("Key(%q): expected ErrInvalidKey, got %v", k, err)
		}
	}
}

func TestKeyStrictReserved(t *testing.T) {
	reservedNames := []string{"PATH", "HOME", "USER", "SHELL", "TERM", "path", "home"}
	for _, k := range reservedNames {
		err := validate.KeyStrict(k)
		if !errors.Is(err, validate.ErrReservedKey) {
			t.Errorf("KeyStrict(%q): expected ErrReservedKey, got %v", k, err)
		}
	}
}

func TestKeyStrictAllowsNonReserved(t *testing.T) {
	if err := validate.KeyStrict("MY_TOKEN"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeysAllValid(t *testing.T) {
	if err := validate.Keys([]string{"FOO", "BAR", "_BAZ"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeysFirstInvalid(t *testing.T) {
	err := validate.Keys([]string{"GOOD", "1BAD", "ALSO_GOOD"})
	if !errors.Is(err, validate.ErrInvalidKey) {
		t.Fatalf("expected ErrInvalidKey, got %v", err)
	}
}
