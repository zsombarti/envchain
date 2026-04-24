package audit_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"envchain/internal/audit"
)

func fixedLogger(buf *bytes.Buffer) *audit.Logger {
	l := audit.New(buf)
	return l
}

func TestLogGet(t *testing.T) {
	var buf bytes.Buffer
	l := fixedLogger(&buf)
	e := l.Get("MY_KEY", "profile:dev")

	if e.Kind != audit.EventGet {
		t.Fatalf("expected GET, got %s", e.Kind)
	}
	if e.Key != "MY_KEY" {
		t.Fatalf("expected key MY_KEY, got %s", e.Key)
	}
	if e.Source != "profile:dev" {
		t.Fatalf("expected source profile:dev, got %s", e.Source)
	}
	if !strings.Contains(buf.String(), "[GET]") {
		t.Errorf("output missing [GET]: %s", buf.String())
	}
}

func TestLogSet(t *testing.T) {
	var buf bytes.Buffer
	l := fixedLogger(&buf)
	e := l.Set("DB_URL", "memory")

	if e.Kind != audit.EventSet {
		t.Fatalf("expected SET, got %s", e.Kind)
	}
	if !strings.Contains(buf.String(), "[SET]") {
		t.Errorf("output missing [SET]: %s", buf.String())
	}
}

func TestLogDelete(t *testing.T) {
	var buf bytes.Buffer
	l := fixedLogger(&buf)
	e := l.Delete("OLD_KEY", "env")

	if e.Kind != audit.EventDelete {
		t.Fatalf("expected DELETE, got %s", e.Kind)
	}
}

func TestLogResolve(t *testing.T) {
	var buf bytes.Buffer
	l := fixedLogger(&buf)
	e := l.Resolve("API_TOKEN", "chain:default")

	if e.Kind != audit.EventResolve {
		t.Fatalf("expected RESOLVE, got %s", e.Kind)
	}
	if !strings.Contains(buf.String(), "chain:default") {
		t.Errorf("output missing source: %s", buf.String())
	}
}

func TestEventTimestamp(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	before := time.Now()
	e := l.Get("X", "src")
	after := time.Now()

	if e.Time.Before(before) || e.Time.After(after) {
		t.Errorf("event time %v out of range [%v, %v]", e.Time, before, after)
	}
}

func TestNilWriterDiscards(t *testing.T) {
	l := audit.New(nil)
	// should not panic
	l.Get("KEY", "src")
	l.Set("KEY", "src")
	l.Delete("KEY", "src")
}
