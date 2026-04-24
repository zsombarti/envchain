package export_test

import (
	"strings"
	"testing"

	"envchain/internal/export"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_URL":   "postgres://localhost/mydb",
	"SECRET":   "p@ss w0rd!",
	"SIMPLE":   "hello",
}

func TestWritePosix(t *testing.T) {
	e := export.New(export.FormatPosix)
	var buf strings.Builder
	if err := e.Write(sampleEnv, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=production\n") {
		t.Errorf("expected posix export line, got:\n%s", out)
	}
	if !strings.Contains(out, "export SECRET=") {
		t.Errorf("expected SECRET to be present, got:\n%s", out)
	}
	// Values with spaces/special chars should be quoted
	if !strings.Contains(out, "'") {
		t.Errorf("expected single-quoted value for SECRET, got:\n%s", out)
	}
}

func TestWriteDotenv(t *testing.T) {
	e := export.New(export.FormatDotenv)
	var buf strings.Builder
	if err := e.Write(sampleEnv, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "export ") {
		t.Errorf("dotenv format must not contain 'export', got:\n%s", out)
	}
	if !strings.Contains(out, "SIMPLE=hello\n") {
		t.Errorf("expected SIMPLE=hello line, got:\n%s", out)
	}
}

func TestWriteJSON(t *testing.T) {
	e := export.New(export.FormatJSON)
	var buf strings.Builder
	if err := e.Write(sampleEnv, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "{\n") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV": "production"`) {
		t.Errorf("expected APP_ENV in JSON, got:\n%s", out)
	}
}

func TestWriteUnknownFormat(t *testing.T) {
	e := export.New(export.Format("xml"))
	var buf strings.Builder
	if err := e.Write(sampleEnv, &buf); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestWriteEmptyEnv(t *testing.T) {
	for _, f := range []export.Format{export.FormatPosix, export.FormatDotenv, export.FormatJSON} {
		e := export.New(f)
		var buf strings.Builder
		if err := e.Write(map[string]string{}, &buf); err != nil {
			t.Errorf("format %s: unexpected error: %v", f, err)
		}
	}
}
