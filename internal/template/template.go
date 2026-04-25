// Package template provides variable interpolation for environment values,
// allowing references like ${VAR} or $VAR to be expanded using a provided map.
package template

import (
	"fmt"
	"strings"
	"unicode"
)

// Options controls expansion behaviour.
type Options struct {
	// AllowMissing suppresses errors for undefined references, substituting
	// an empty string instead.
	AllowMissing bool

	// MaxDepth limits recursive expansion passes to prevent cycles.
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		AllowMissing: false,
		MaxDepth:     10,
	}
}

// Expand replaces ${VAR} and $VAR references in s using the values in env.
// It returns an error if a referenced key is missing and AllowMissing is false.
func Expand(s string, env map[string]string, opts Options) (string, error) {
	prev := s
	for i := 0; i < opts.MaxDepth; i++ {
		result, err := expandOnce(prev, env, opts)
		if err != nil {
			return "", err
		}
		if result == prev {
			return result, nil
		}
		prev = result
	}
	return prev, nil
}

// ExpandMap applies Expand to every value in m, returning a new map.
func ExpandMap(m map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		expanded, err := Expand(v, m, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}

func expandOnce(s string, env map[string]string, opts Options) (string, error) {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '$' {
			b.WriteByte(s[i])
			i++
			continue
		}
		// escape: $$
		if i+1 < len(s) && s[i+1] == '$' {
			b.WriteByte('$')
			i += 2
			continue
		}
		var key string
		var advance int
		if i+1 < len(s) && s[i+1] == '{' {
			end := strings.IndexByte(s[i+2:], '}')
			if end < 0 {
				return "", fmt.Errorf("unclosed '${' in %q", s)
			}
			key = s[i+2 : i+2+end]
			advance = 3 + end
		} else {
			j := i + 1
			for j < len(s) && isIdentRune(rune(s[j])) {
				j++
			}
			key = s[i+1 : j]
			advance = j - i
		}
		if key == "" {
			b.WriteByte('$')
			i++
			continue
		}
		val, ok := env[key]
		if !ok {
			if !opts.AllowMissing {
				return "", fmt.Errorf("undefined variable %q", key)
			}
		}
		b.WriteString(val)
		i += advance
	}
	return b.String(), nil
}

func isIdentRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
