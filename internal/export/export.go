// Package export provides utilities for rendering resolved environment
// variable chains into various shell-compatible output formats.
package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	// FormatPosix emits POSIX-compatible export statements (export KEY=VALUE).
	FormatPosix Format = "posix"
	// FormatDotenv emits dotenv-compatible KEY=VALUE pairs.
	FormatDotenv Format = "dotenv"
	// FormatJSON emits a JSON object of key/value pairs.
	FormatJSON Format = "json"
)

// Exporter writes a resolved environment map to an io.Writer in the
// requested format.
type Exporter struct {
	format Format
}

// New creates a new Exporter for the given Format.
func New(format Format) *Exporter {
	return &Exporter{format: format}
}

// Write renders the environment map to w.
func (e *Exporter) Write(env map[string]string, w io.Writer) error {
	keys := sortedKeys(env)
	switch e.format {
	case FormatPosix:
		return writePosix(env, keys, w)
	case FormatDotenv:
		return writeDotenv(env, keys, w)
	case FormatJSON:
		return writeJSON(env, keys, w)
	default:
		return fmt.Errorf("export: unknown format %q", e.format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writePosix(env map[string]string, keys []string, w io.Writer) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, shellescape(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeDotenv(env map[string]string, keys []string, w io.Writer) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, shellescape(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(env map[string]string, keys []string, w io.Writer) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		sb.WriteString(fmt.Sprintf("  %q: %q", k, env[k]))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// shellescape wraps a value in single quotes if it contains special characters.
func shellescape(v string) string {
	if strings.ContainsAny(v, " \t\n\r'\"\\$`!#&;|<>(){}") {
		return "'" + strings.ReplaceAll(v, "'", "'\\''")+"'"
	}
	return v
}
