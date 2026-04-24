package store

import (
	"errors"
)

// ErrSecretNotFound is returned when a secret key does not exist in the store.
var ErrSecretNotFound = errors.New("secret not found")

// Store defines the interface for a secret backend.
type Store interface {
	// Get retrieves a secret value by key.
	Get(key string) (string, error)
	// Set stores a secret value by key.
	Set(key, value string) error
	// Delete removes a secret by key.
	Delete(key string) error
	// List returns all secret keys managed by this store.
	List() ([]string, error)
}
