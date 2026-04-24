package redact_test

import (
	"testing"

	"envchain/internal/redact"
)

func TestValueDefault(t *testing.T) {
	got := redact.Value("supersecret", redact.DefaultOptions())
	if got != redact.DefaultMask {
		t.Fatalf("expected %q, got %q", redact.DefaultMask, got)
	}
}

func TestValuePartial(t *testing.T) {
	opts := redact.Options{Mask: "****", Partial: true}
	got := redact.Value("supersecret", opts)
	want := "****" + "cret"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestValuePartialShort(t *testing.T) {
	// Value shorter than PartialReveal should still be fully masked.
	opts := redact.Options{Mask: "****", Partial: true}
	got := redact.Value("abc", opts)
	if got != "****" {
		t.Fatalf("expected full mask for short value, got %q", got)
	}
}

func TestValueEmptyMaskFallback(t *testing.T) {
	opts := redact.Options{Mask: "", Partial: false}
	got := redact.Value("myvalue", opts)
	if got != redact.DefaultMask {
		t.Fatalf("expected default mask fallback, got %q", got)
	}
}

func TestMapRedactsSecrets(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "abc123",
		"USERNAME": "alice",
		"PASSWORD": "s3cr3t",
	}
	secrets := []string{"API_KEY", "PASSWORD"}
	out := redact.Map(env, secrets, redact.DefaultOptions())

	if out["API_KEY"] != redact.DefaultMask {
		t.Errorf("API_KEY should be redacted, got %q", out["API_KEY"])
	}
	if out["PASSWORD"] != redact.DefaultMask {
		t.Errorf("PASSWORD should be redacted, got %q", out["PASSWORD"])
	}
	if out["USERNAME"] != "alice" {
		t.Errorf("USERNAME should be unchanged, got %q", out["USERNAME"])
	}
}

func TestMapCaseInsensitiveKeys(t *testing.T) {
	env := map[string]string{"Api_Key": "secret"}
	out := redact.Map(env, []string{"api_key"}, redact.DefaultOptions())
	if out["Api_Key"] != redact.DefaultMask {
		t.Errorf("case-insensitive match failed, got %q", out["Api_Key"])
	}
}

func TestMapEmptySecrets(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	out := redact.Map(env, nil, redact.DefaultOptions())
	if out["FOO"] != "bar" {
		t.Errorf("expected value unchanged, got %q", out["FOO"])
	}
}

func TestKeysFindsMatches(t *testing.T) {
	env := map[string]string{"TOKEN": "x", "HOST": "localhost", "PORT": "8080"}
	found := redact.Keys(env, []string{"TOKEN", "PORT"})
	if len(found) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(found), found)
	}
}

func TestKeysNoMatch(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	found := redact.Keys(env, []string{"MISSING"})
	if len(found) != 0 {
		t.Fatalf("expected no keys, got %v", found)
	}
}
