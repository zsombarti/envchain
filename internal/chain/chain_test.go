package chain

import (
	"testing"
)

func TestAddLayer(t *testing.T) {
	c := New()
	err := c.AddLayer("base", map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Layers) != 1 {
		t.Fatalf("expected 1 layer, got %d", len(c.Layers))
	}
}

func TestAddLayerDuplicate(t *testing.T) {
	c := New()
	_ = c.AddLayer("base", map[string]string{"FOO": "bar"})
	err := c.AddLayer("base", map[string]string{"BAZ": "qux"})
	if err == nil {
		t.Fatal("expected error for duplicate layer name")
	}
}

func TestRemoveLayer(t *testing.T) {
	c := New()
	_ = c.AddLayer("base", map[string]string{"FOO": "bar"})
	err := c.RemoveLayer("base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Layers) != 0 {
		t.Fatalf("expected 0 layers, got %d", len(c.Layers))
	}
}

func TestRemoveLayerNotFound(t *testing.T) {
	c := New()
	err := c.RemoveLayer("missing")
	if err == nil {
		t.Fatal("expected error for missing layer")
	}
}

func TestResolveEmpty(t *testing.T) {
	c := New()
	_, err := c.Resolve()
	if err == nil {
		t.Fatal("expected error when resolving empty chain")
	}
}

func TestResolvePrecedence(t *testing.T) {
	c := New()
	_ = c.AddLayer("base", map[string]string{"FOO": "base", "BAR": "base"})
	_ = c.AddLayer("override", map[string]string{"FOO": "override", "BAZ": "new"})

	env, err := c.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "override" {
		t.Errorf("expected FOO=override, got %s", env["FOO"])
	}
	if env["BAR"] != "base" {
		t.Errorf("expected BAR=base, got %s", env["BAR"])
	}
	if env["BAZ"] != "new" {
		t.Errorf("expected BAZ=new, got %s", env["BAZ"])
	}
}

func TestGetLayer(t *testing.T) {
	c := New()
	_ = c.AddLayer("base", map[string]string{"FOO": "bar"})
	l, err := c.GetLayer("base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Name != "base" {
		t.Errorf("expected layer name 'base', got %s", l.Name)
	}
}
