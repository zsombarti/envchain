// Package resolve provides utilities for resolving environment variables
// from a chain of layers, applying overrides and filtering by prefix.
package resolve

import (
	"strings"

	"envchain/internal/chain"
)

// Options configures the resolution behaviour.
type Options struct {
	// Prefix filters keys to only those starting with the given string.
	// The prefix is stripped from the resulting keys when StripPrefix is true.
	Prefix string
	// StripPrefix removes the Prefix from resolved key names.
	StripPrefix bool
	// Overrides are applied last, taking highest precedence.
	Overrides map[string]string
}

// DefaultOptions returns an Options with no filtering or overrides.
func DefaultOptions() Options {
	return Options{}
}

// Resolve walks the chain and returns a merged map of environment variables
// according to the provided Options.
func Resolve(c *chain.Chain, opts Options) (map[string]string, error) {
	raw, err := c.Resolve()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(raw))

	for k, v := range raw {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		key := k
		if opts.StripPrefix && opts.Prefix != "" {
			key = strings.TrimPrefix(k, opts.Prefix)
		}
		result[key] = v
	}

	for k, v := range opts.Overrides {
		result[k] = v
	}

	return result, nil
}

// Keys returns the sorted list of resolved keys.
func Keys(c *chain.Chain, opts Options) ([]string, error) {
	env, err := Resolve(c, opts)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys, nil
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
