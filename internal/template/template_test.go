package template

import (
	"testing"
)

func TestExpandSimple(t *testing.T) {
	env := map[string]string{"HOME": "/home/user"}
	got, err := Expand("path=$HOME/bin", env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "path=/home/user/bin" {
		t.Fatalf("got %q", got)
	}
}

func TestExpandBraces(t *testing.T) {
	env := map[string]string{"USER": "alice"}
	got, err := Expand("hello_${USER}_world", env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello_alice_world" {
		t.Fatalf("got %q", got)
	}
}

func TestExpandMissingError(t *testing.T) {
	env := map[string]string{}
	_, err := Expand("$MISSING", env, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestExpandMissingAllowed(t *testing.T) {
	env := map[string]string{}
	opts := DefaultOptions()
	opts.AllowMissing = true
	got, err := Expand("prefix_${MISSING}_suffix", env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "prefix__suffix" {
		t.Fatalf("got %q", got)
	}
}

func TestExpandEscapeDollar(t *testing.T) {
	env := map[string]string{"A": "1"}
	got, err := Expand("cost=$$5", env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "cost=$5" {
		t.Fatalf("got %q", got)
	}
}

func TestExpandUnclosedBrace(t *testing.T) {
	env := map[string]string{}
	_, err := Expand("${UNCLOSED", env, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for unclosed brace")
	}
}

func TestExpandRecursive(t *testing.T) {
	env := map[string]string{
		"BASE": "/usr",
		"BIN":  "${BASE}/bin",
	}
	got, err := Expand("$BIN", env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/usr/bin" {
		t.Fatalf("got %q", got)
	}
}

func TestExpandMapRoundtrip(t *testing.T) {
	m := map[string]string{
		"PREFIX": "pg",
		"HOST":   "${PREFIX}-host",
		"DSN":    "postgres://${HOST}:5432",
	}
	out, err := ExpandMap(m, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://pg-host:5432" {
		t.Fatalf("got DSN=%q", out["DSN"])
	}
}

func TestExpandMapMissingError(t *testing.T) {
	m := map[string]string{
		"VAL": "$UNDEFINED",
	}
	_, err := ExpandMap(m, DefaultOptions())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExpandNoOp(t *testing.T) {
	env := map[string]string{}
	got, err := Expand("no refs here", env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "no refs here" {
		t.Fatalf("got %q", got)
	}
}
