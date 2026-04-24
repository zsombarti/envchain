package prompt

import (
	"bytes"
	"os"
	"testing"
)

// pipeFile creates an *os.File backed by a pipe whose read end contains data.
func pipeFile(t *testing.T, data string) *os.File {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	if _, err := w.WriteString(data); err != nil {
		t.Fatalf("write pipe: %v", err)
	}
	w.Close()
	t.Cleanup(func() { r.Close() })
	return r
}

func TestPassphraseNoConfirm(t *testing.T) {
	in := pipeFile(t, "s3cr3t\n")
	out := &bytes.Buffer{}
	opts := &Options{In: in, Out: out}

	got, err := Passphrase("Passphrase", false, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "s3cr3t" {
		t.Errorf("got %q, want %q", got, "s3cr3t")
	}
}

func TestPassphraseConfirmMatch(t *testing.T) {
	in := pipeFile(t, "s3cr3t\ns3cr3t\n")
	out := &bytes.Buffer{}
	opts := &Options{In: in, Out: out}

	got, err := Passphrase("Passphrase", true, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "s3cr3t" {
		t.Errorf("got %q, want %q", got, "s3cr3t")
	}
}

func TestPassphraseConfirmMismatch(t *testing.T) {
	in := pipeFile(t, "s3cr3t\ndifferent\n")
	out := &bytes.Buffer{}
	opts := &Options{In: in, Out: out}

	_, err := Passphrase("Passphrase", true, opts)
	if err != ErrMismatch {
		t.Errorf("expected ErrMismatch, got %v", err)
	}
}

func TestConfirmYes(t *testing.T) {
	for _, input := range []string{"y\n", "Y\n", "yes\n", "YES\n"} {
		in := pipeFile(t, input)
		out := &bytes.Buffer{}
		opts := &Options{In: in, Out: out}

		ok, err := Confirm("Continue?", opts)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if !ok {
			t.Errorf("input %q: expected true", input)
		}
	}
}

func TestConfirmNo(t *testing.T) {
	for _, input := range []string{"n\n", "N\n", "no\n", "\n"} {
		in := pipeFile(t, input)
		out := &bytes.Buffer{}
		opts := &Options{In: in, Out: out}

		ok, err := Confirm("Continue?", opts)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if ok {
			t.Errorf("input %q: expected false", input)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	o := defaults(nil)
	if o.In == nil {
		t.Error("expected In to default to os.Stdin")
	}
	if o.Out == nil {
		t.Error("expected Out to default to os.Stderr")
	}
}
