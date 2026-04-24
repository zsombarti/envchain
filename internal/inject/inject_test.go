package inject_test

import (
	"testing"

	"envchain/internal/inject"
)

func TestInjectBasic(t *testing.T) {
	dst := map[string]string{"EXISTING": "old"}
	src := map[string]string{"NEW": "value"}
	skipped, err := inject.Inject(dst, src, inject.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", skipped)
	}
	if dst["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %q", dst["NEW"])
	}
}

func TestInjectNoOverwrite(t *testing.T) {
	dst := map[string]string{"KEY": "original"}
	src := map[string]string{"KEY": "new"}
	skipped, err := inject.Inject(dst, src, inject.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := skipped["KEY"]; !ok {
		t.Error("expected KEY to be skipped")
	}
	if dst["KEY"] != "original" {
		t.Errorf("expected KEY=original, got %q", dst["KEY"])
	}
}

func TestInjectOverwrite(t *testing.T) {
	dst := map[string]string{"KEY": "original"}
	src := map[string]string{"KEY": "new"}
	opts := inject.Options{Overwrite: true}
	_, err := inject.Inject(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", dst["KEY"])
	}
}

func TestInjectWithPrefix(t *testing.T) {
	dst := make(map[string]string)
	src := map[string]string{"VAR": "val"}
	opts := inject.Options{Prefix: "APP_"}
	_, err := inject.Inject(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["APP_VAR"] != "val" {
		t.Errorf("expected APP_VAR=val, got %q", dst["APP_VAR"])
	}
}

func TestInjectNilDst(t *testing.T) {
	_, err := inject.Inject(nil, map[string]string{}, inject.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMerge(t *testing.T) {
	a := map[string]string{"A": "1", "SHARED": "from-a"}
	b := map[string]string{"B": "2", "SHARED": "from-b"}
	opts := inject.Options{Overwrite: true}
	result, err := inject.Merge(opts, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %q", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2, got %q", result["B"])
	}
	if result["SHARED"] != "from-b" {
		t.Errorf("expected SHARED=from-b, got %q", result["SHARED"])
	}
}
