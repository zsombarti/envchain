// Package inject provides utilities for merging and injecting environment
// variable sets into a target map, respecting precedence and override rules.
package inject

import "fmt"

// Options controls how injection behaves.
type Options struct {
	// Overwrite allows injected values to replace existing keys in the target.
	Overwrite bool
	// Prefix is prepended to every injected key.
	Prefix string
}

// DefaultOptions returns Options with safe defaults: no overwrite, no prefix.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Prefix:    "",
	}
}

// Inject merges src into dst according to opts.
// It returns a map of keys that were skipped because they already existed
// in dst and Overwrite was false.
func Inject(dst, src map[string]string, opts Options) (skipped map[string]string, err error) {
	if dst == nil {
		return nil, fmt.Errorf("inject: dst must not be nil")
	}
	skipped = make(map[string]string)
	for k, v := range src {
		key := opts.Prefix + k
		if _, exists := dst[key]; exists && !opts.Overwrite {
			skipped[key] = v
			continue
		}
		dst[key] = v
	}
	return skipped, nil
}

// MustInject is like Inject but panics on error.
func MustInject(dst, src map[string]string, opts Options) map[string]string {
	skipped, err := Inject(dst, src, opts)
	if err != nil {
		panic(err)
	}
	return skipped
}

// Merge combines multiple source maps into a new map using the provided opts.
// Sources are applied in order; later sources win when Overwrite is true.
func Merge(opts Options, sources ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	for _, src := range sources {
		_, err := Inject(result, src, opts)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
