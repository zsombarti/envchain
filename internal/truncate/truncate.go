// Package truncate provides utilities for truncating long strings in
// CLI output, ensuring values fit within a given display width.
package truncate

const (
	DefaultMaxLen  = 64
	DefaultSuffix  = "..."
	MinDisplayLen  = 8
)

// Options controls truncation behaviour.
type Options struct {
	// MaxLen is the maximum number of runes to display (including suffix).
	// Defaults to DefaultMaxLen if zero.
	MaxLen int
	// Suffix is appended when a value is truncated. Defaults to DefaultSuffix.
	Suffix string
}

func defaults(o Options) Options {
	if o.MaxLen <= 0 {
		o.MaxLen = DefaultMaxLen
	}
	if o.Suffix == "" {
		o.Suffix = DefaultSuffix
	}
	return o
}

// Value truncates s to at most opts.MaxLen runes. If s is longer than
// MaxLen, the returned string is shortened and opts.Suffix is appended.
// The total length of the returned string never exceeds MaxLen runes.
func Value(s string, opts Options) string {
	opts = defaults(opts)
	runes := []rune(s)
	if len(runes) <= opts.MaxLen {
		return s
	}
	suffixRunes := []rune(opts.Suffix)
	cutAt := opts.MaxLen - len(suffixRunes)
	if cutAt < 0 {
		cutAt = 0
	}
	return string(runes[:cutAt]) + opts.Suffix
}

// Map applies Value to every entry in m, returning a new map with the
// same keys and truncated values.
func Map(m map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = Value(v, opts)
	}
	return out
}

// Key truncates an environment variable key using a stricter default
// max length suited for columnar display.
func Key(s string, maxLen int) string {
	if maxLen <= 0 {
		maxLen = 32
	}
	return Value(s, Options{MaxLen: maxLen, Suffix: DefaultSuffix})
}
