package resolve_test

import (
	"strings"
	"testing"

	"envchain/internal/resolve"
)

func TestMergeEmpty(t *testing.T) {
	out := resolve.Merge()
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestMergePrecedence(t *testing.T) {
	a := map[string]string{"X": "first", "Y": "only"}
	b := map[string]string{"X": "second", "Z": "new"}
	out := resolve.Merge(a, b)
	if out["X"] != "second" {
		t.Fatalf("expected second, got %q", out["X"])
	}
	if out["Y"] != "only" {
		t.Fatalf("expected only, got %q", out["Y"])
	}
	if out["Z"] != "new" {
		t.Fatalf("expected new, got %q", out["Z"])
	}
}

func TestFilter(t *testing.T) {
	env := map[string]string{"KEEP_A": "1", "KEEP_B": "2", "DROP": "3"}
	out := resolve.Filter(env, func(k, _ string) bool {
		return strings.HasPrefix(k, "KEEP_")
	})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["DROP"]; ok {
		t.Fatal("DROP should have been filtered")
	}
}

func TestRename(t *testing.T) {
	env := map[string]string{"OLD_FOO": "bar", "OLD_BAZ": "qux"}
	out := resolve.Rename(env, func(k string) string {
		return strings.TrimPrefix(k, "OLD_")
	})
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected rename result: %v", out)
	}
}

func TestRenameDropEmpty(t *testing.T) {
	env := map[string]string{"KEEP": "1", "DROP": "2"}
	out := resolve.Rename(env, func(k string) string {
		if k == "DROP" {
			return ""
		}
		return k
	})
	if len(out) != 1 || out["KEEP"] != "1" {
		t.Fatalf("unexpected result: %v", out)
	}
}
