package resolve_test

import (
	"fmt"
	"testing"

	"envchain/internal/chain"
	"envchain/internal/resolve"
	"envchain/internal/store"
)

func buildChain(t *testing.T, layers ...map[string]string) *chain.Chain {
	t.Helper()
	c := chain.New()
	for i, kv := range layers {
		name := fmt.Sprintf("layer%d", i)
		ms := store.NewMemoryStore()
		for k, v := range kv {
			_ = ms.Set(k, v)
		}
		if err := c.AddLayer(name, ms); err != nil {
			t.Fatalf("AddLayer: %v", err)
		}
	}
	return c
}

func TestResolveBasic(t *testing.T) {
	c := buildChain(t, map[string]string{"FOO": "bar", "BAZ": "qux"})
	env, err := resolve.Resolve(c, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" || env["BAZ"] != "qux" {
		t.Fatalf("unexpected env: %v", env)
	}
}

func TestResolvePrefix(t *testing.T) {
	c := buildChain(t, map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"})
	env, err := resolve.Resolve(c, resolve.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(env), env)
	}
	if _, ok := env["OTHER"]; ok {
		t.Fatal("OTHER should have been filtered out")
	}
}

func TestResolveStripPrefix(t *testing.T) {
	c := buildChain(t, map[string]string{"APP_FOO": "1", "APP_BAR": "2"})
	env, err := resolve.Resolve(c, resolve.Options{Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "1" || env["BAR"] != "2" {
		t.Fatalf("unexpected env after strip: %v", env)
	}
}

func TestResolveOverrides(t *testing.T) {
	c := buildChain(t, map[string]string{"FOO": "original"})
	env, err := resolve.Resolve(c, resolve.Options{Overrides: map[string]string{"FOO": "override", "NEW": "val"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "override" {
		t.Fatalf("expected override, got %q", env["FOO"])
	}
	if env["NEW"] != "val" {
		t.Fatalf("expected NEW=val, got %q", env["NEW"])
	}
}

func TestKeys(t *testing.T) {
	c := buildChain(t, map[string]string{"Z": "1", "A": "2", "M": "3"})
	keys, err := resolve.Keys(c, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("unexpected keys order: %v", keys)
	}
}
