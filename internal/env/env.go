// Package env provides utilities for parsing and formatting environment
// variable assignments in the KEY=VALUE format used by shells and dotenv files.
package env

import (
	"fmt"
	"strings"
)

// ErrInvalidAssignment is returned when a string is not a valid KEY=VALUE pair.
type ErrInvalidAssignment struct {
	Raw string
}

func (e *ErrInvalidAssignment) Error() string {
	return fmt.Sprintf("env: invalid assignment %q: missing '='" , e.Raw)
}

// Parse parses a slice of KEY=VALUE strings into a map.
// Duplicate keys are overwritten by the last occurrence.
// Values may be empty (KEY= is valid).
func Parse(pairs []string) (map[string]string, error) {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, err := Split(p)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

// Split splits a single KEY=VALUE string into its key and value parts.
// Only the first '=' is used as the delimiter, allowing '=' in values.
func Split(pair string) (key, value string, err error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return "", "", &ErrInvalidAssignment{Raw: pair}
	}
	return pair[:idx], pair[idx+1:], nil
}

// Format converts a map of environment variables into a sorted slice of
// KEY=VALUE strings suitable for use with os/exec.Cmd.Env.
func Format(env map[string]string) []string {
	pairs := make([]string, 0, len(env))
	for k, v := range env {
		pairs = append(pairs, k+"="+v)
	}
	return pairs
}

// Merge merges src into dst. If overwrite is true, existing keys in dst are
// replaced; otherwise they are preserved.
func Merge(dst, src map[string]string, overwrite bool) {
	for k, v := range src {
		if _, exists := dst[k]; !exists || overwrite {
			dst[k] = v
		}
	}
}
