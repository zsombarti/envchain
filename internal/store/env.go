package store

import (
	"fmt"
	"os"
	"strings"
)

// EnvStore reads secrets from the current process environment.
// It is read-only; Set and Delete are no-ops that return errors.
type EnvStore struct {
	prefix string
}

// NewEnvStore creates an EnvStore that optionally filters by a key prefix.
// Pass an empty string to match all environment variables.
func NewEnvStore(prefix string) *EnvStore {
	return &EnvStore{prefix: prefix}
}

// Get retrieves an environment variable by key (prefix is not prepended).
func (e *EnvStore) Get(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", ErrSecretNotFound
	}
	return v, nil
}

// Set is not supported for the environment store.
func (e *EnvStore) Set(_, _ string) error {
	return fmt.Errorf("EnvStore is read-only")
}

// Delete is not supported for the environment store.
func (e *EnvStore) Delete(_ string) error {
	return fmt.Errorf("EnvStore is read-only")
}

// List returns all environment variable keys that match the configured prefix.
func (e *EnvStore) List() ([]string, error) {
	env := os.Environ()
	keys := make([]string, 0, len(env))
	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) < 1 {
			continue
		}
		key := parts[0]
		if e.prefix == "" || strings.HasPrefix(key, e.prefix) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}
