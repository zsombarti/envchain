package truncate_test

import (
	"strings"
	"testing"

	"envchain/internal/truncate"
)

func TestValueShortString(t *testing.T) {
	s := "SHORT"
	got := truncate.Value(s, truncate.Options{MaxLen: 20})
	if got != s {
		t.Fatalf("expected %q, got %q", s, got)
	}
}

func TestValueExactLength(t *testing.T) {
	s := strings.Repeat("x", 64)
	got := truncate.Value(s, truncate.Options{})
	if got != s {
		t.Fatalf("expected unchanged string of length 64")
	}
}

func TestValueTruncated(t *testing.T) {
	s := strings.Repeat("a", 80)
	got := truncate.Value(s, truncate.Options{MaxLen: 10, Suffix: "..."})
	if len([]rune(got)) != 10 {
		t.Fatalf("expected length 10, got %d", len([]rune(got)))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestValueCustomSuffix(t *testing.T) {
	s := "hello world this is long"
	got := truncate.Value(s, truncate.Options{MaxLen: 10, Suffix: "~"})
	if !strings.HasSuffix(got, "~") {
		t.Fatalf("expected '~' suffix, got %q", got)
	}
	if len([]rune(got)) > 10 {
		t.Fatalf("result exceeds MaxLen: %q", got)
	}
}

func TestValueUnicode(t *testing.T) {
	// Each character is a multi-byte rune.
	s := strings.Repeat("é", 20)
	got := truncate.Value(s, truncate.Options{MaxLen: 10, Suffix: "..."})
	if len([]rune(got)) > 10 {
		t.Fatalf("unicode truncation exceeded MaxLen: %q", got)
	}
}

func TestMapTruncatesValues(t *testing.T) {
	m := map[string]string{
		"SHORT": "hi",
		"LONG":  strings.Repeat("z", 100),
	}
	out := truncate.Map(m, truncate.Options{MaxLen: 20})
	if out["SHORT"] != "hi" {
		t.Fatalf("short value should be unchanged")
	}
	if len([]rune(out["LONG"])) > 20 {
		t.Fatalf("long value should be truncated")
	}
}

func TestMapDoesNotMutateInput(t *testing.T) {
	orig := strings.Repeat("x", 80)
	m := map[string]string{"K": orig}
	truncate.Map(m, truncate.Options{MaxLen: 10})
	if m["K"] != orig {
		t.Fatal("input map was mutated")
	}
}

func TestKey(t *testing.T) {
	long := strings.Repeat("A", 50)
	got := truncate.Key(long, 10)
	if len([]rune(got)) > 10 {
		t.Fatalf("key exceeds maxLen: %q", got)
	}
}

func TestKeyDefaultMaxLen(t *testing.T) {
	short := "MY_VAR"
	got := truncate.Key(short, 0)
	if got != short {
		t.Fatalf("short key should be unchanged, got %q", got)
	}
}
