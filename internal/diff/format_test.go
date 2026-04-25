package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"envchain/internal/diff"
)

func TestWriteAdded(t *testing.T) {
	r := diff.Compute(map[string]string{}, map[string]string{"FOO": "bar"})
	var buf bytes.Buffer
	if err := diff.Write(&buf, r, diff.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "+ FOO=bar") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestWriteRemoved(t *testing.T) {
	r := diff.Compute(map[string]string{"FOO": "bar"}, map[string]string{})
	var buf bytes.Buffer
	if err := diff.Write(&buf, r, diff.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "- FOO=bar") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestWriteModified(t *testing.T) {
	r := diff.Compute(map[string]string{"X": "old"}, map[string]string{"X": "new"})
	var buf bytes.Buffer
	if err := diff.Write(&buf, r, diff.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "~ X: old -> new") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestWriteRedacted(t *testing.T) {
	r := diff.Compute(map[string]string{}, map[string]string{"SECRET": "s3cr3t"})
	var buf bytes.Buffer
	opts := diff.Options{RedactValues: true}
	if err := diff.Write(&buf, r, opts); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "s3cr3t") {
		t.Fatal("value should have been redacted")
	}
	if !strings.Contains(buf.String(), "***") {
		t.Fatalf("expected redaction marker, got: %q", buf.String())
	}
}

func TestWriteUnchangedHidden(t *testing.T) {
	r := diff.Compute(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	var buf bytes.Buffer
	if err := diff.Write(&buf, r, diff.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output for unchanged, got: %q", buf.String())
	}
}

func TestSummary(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"B": "changed", "C": "3"}
	r := diff.Compute(old, next)
	s := diff.Summary(r)
	if !strings.Contains(s, "+1") || !strings.Contains(s, "-1") || !strings.Contains(s, "~1") {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestSummaryNoChanges(t *testing.T) {
	r := diff.Compute(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	if diff.Summary(r) != "no changes" {
		t.Fatalf("expected 'no changes', got: %q", diff.Summary(r))
	}
}
