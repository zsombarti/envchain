package env_test

import (
	"testing"

	"envchain/internal/env"
)

func TestParse(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux", "EMPTY="}
	got, err := env.Parse(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("FOO: want bar, got %q", got["FOO"])
	}
	if got["BAZ"] != "qux" {
		t.Errorf("BAZ: want qux, got %q", got["BAZ"])
	}
	if got["EMPTY"] != "" {
		t.Errorf("EMPTY: want empty string, got %q", got["EMPTY"])
	}
}

func TestParseInvalid(t *testing.T) {
	_, err := env.Parse([]string{"NOEQUALS"})
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
	var e *env.ErrInvalidAssignment
	if _, ok := err.(*env.ErrInvalidAssignment); !ok {
		t.Errorf("expected *ErrInvalidAssignment, got %T", e)
	}
}

func TestSplitValueContainsEquals(t *testing.T) {
	k, v, err := env.Split("URL=http://example.com?a=1&b=2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k != "URL" {
		t.Errorf("key: want URL, got %q", k)
	}
	if v != "http://example.com?a=1&b=2" {
		t.Errorf("value: want full URL, got %q", v)
	}
}

func TestFormatRoundtrip(t *testing.T) {
	input := map[string]string{"A": "1", "B": "2"}
	pairs := env.Format(input)
	got, err := env.Parse(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range input {
		if got[k] != v {
			t.Errorf("%s: want %q, got %q", k, v, got[k])
		}
	}
}

func TestMergeNoOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new", "B": "added"}
	env.Merge(dst, src, false)
	if dst["A"] != "original" {
		t.Errorf("A should not be overwritten, got %q", dst["A"])
	}
	if dst["B"] != "added" {
		t.Errorf("B should be added, got %q", dst["B"])
	}
}

func TestMergeWithOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new"}
	env.Merge(dst, src, true)
	if dst["A"] != "new" {
		t.Errorf("A should be overwritten, got %q", dst["A"])
	}
}
