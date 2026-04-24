// Package redact provides utilities for masking sensitive environment
// variable values in output, logs, and display contexts.
package redact

import (
	"strings"
)

const (
	// DefaultMask is the default replacement string for redacted values.
	DefaultMask = "********"

	// PartialReveal controls how many characters to reveal at the end
	// of a value when using partial redaction.
	PartialReveal = 4
)

// Options configures redaction behaviour.
type Options struct {
	// Mask is the string used to replace redacted values.
	Mask string
	// Partial reveals the last N characters of the value.
	Partial bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Mask:    DefaultMask,
		Partial: false,
	}
}

// Value redacts a single secret value according to opts.
func Value(v string, opts Options) string {
	if opts.Mask == "" {
		opts.Mask = DefaultMask
	}
	if !opts.Partial || len(v) <= PartialReveal {
		return opts.Mask
	}
	return opts.Mask + v[len(v)-PartialReveal:]
}

// Map returns a copy of env where every key listed in secrets has its
// value replaced by the redacted form.
func Map(env map[string]string, secrets []string, opts Options) map[string]string {
	set := make(map[string]struct{}, len(secrets))
	for _, k := range secrets {
		set[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, secret := set[strings.ToUpper(k)]; secret {
			out[k] = Value(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}

// Keys returns the subset of keys from env that appear in secrets.
func Keys(env map[string]string, secrets []string) []string {
	set := make(map[string]struct{}, len(secrets))
	for _, k := range secrets {
		set[strings.ToUpper(k)] = struct{}{}
	}
	var found []string
	for k := range env {
		if _, ok := set[strings.ToUpper(k)]; ok {
			found = append(found, k)
		}
	}
	return found
}
