// Package validate provides helpers for validating environment variable
// names and values before they are stored or exported.
package validate

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ErrEmptyKey is returned when a key is an empty string.
var ErrEmptyKey = errors.New("validate: key must not be empty")

// ErrInvalidKey is returned when a key contains characters that are not
// permitted in POSIX environment variable names.
var ErrInvalidKey = errors.New("validate: key contains invalid characters")

// ErrReservedKey is returned when a key uses a name that is reserved by the
// shell or the tool itself.
var ErrReservedKey = errors.New("validate: key is reserved")

// reserved holds names that must not be used as environment variable keys.
var reserved = map[string]struct{}{
	"PATH": {},
	"HOME": {},
	"USER": {},
	"SHELL": {},
	"TERM": {},
}

// Key returns nil when name is a valid POSIX environment variable name:
// it must be non-empty, start with a letter or underscore, and contain only
// letters, digits, or underscores.
func Key(name string) error {
	if name == "" {
		return ErrEmptyKey
	}
	for i, r := range name {
		switch {
		case i == 0 && !unicode.IsLetter(r) && r != '_':
			return fmt.Errorf("%w: %q must start with a letter or underscore", ErrInvalidKey, name)
		case i > 0 && !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_':
			return fmt.Errorf("%w: %q has invalid character %q at position %d", ErrInvalidKey, name, r, i)
		}
	}
	return nil
}

// KeyStrict is like Key but additionally rejects names in the reserved set.
func KeyStrict(name string) error {
	if err := Key(name); err != nil {
		return err
	}
	if _, ok := reserved[strings.ToUpper(name)]; ok {
		return fmt.Errorf("%w: %q", ErrReservedKey, name)
	}
	return nil
}

// Keys validates every key in names and returns the first error encountered,
// or nil if all keys are valid.
func Keys(names []string) error {
	for _, n := range names {
		if err := Key(n); err != nil {
			return err
		}
	}
	return nil
}
