package diff

import (
	"fmt"
	"io"
	"strings"
)

// Options controls formatting behaviour.
type Options struct {
	// ShowUnchanged includes unchanged keys in the output.
	ShowUnchanged bool
	// RedactValues replaces values with "***" for secret-safe display.
	RedactValues bool
}

// DefaultOptions returns sensible formatting defaults.
func DefaultOptions() Options {
	return Options{ShowUnchanged: false, RedactValues: false}
}

// Write renders a human-readable diff to w using the given options.
func Write(w io.Writer, r Result, opts Options) error {
	for _, c := range r.Changes {
		if c.Kind == Unchanged && !opts.ShowUnchanged {
			continue
		}
		line, err := formatChange(c, opts.RedactValues)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatChange(c Change, redact bool) (string, error) {
	val := func(v string) string {
		if redact {
			return "***"
		}
		return v
	}
	switch c.Kind {
	case Added:
		return fmt.Sprintf("+ %s=%s", c.Key, val(c.NewValue)), nil
	case Removed:
		return fmt.Sprintf("- %s=%s", c.Key, val(c.OldValue)), nil
	case Modified:
		return fmt.Sprintf("~ %s: %s -> %s", c.Key, val(c.OldValue), val(c.NewValue)), nil
	case Unchanged:
		return fmt.Sprintf("  %s=%s", c.Key, val(c.NewValue)), nil
	default:
		return "", fmt.Errorf("diff: unknown change kind %q", c.Kind)
	}
}

// Summary returns a one-line summary string for the result.
func Summary(r Result) string {
	var added, removed, modified int
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	parts := make([]string, 0, 3)
	if added > 0 {
		parts = append(parts, fmt.Sprintf("+%d", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d", removed))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("~%d", modified))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, " ")
}
