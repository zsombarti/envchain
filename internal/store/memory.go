package store

import (
	"sync"
)

// MemoryStore is an in-memory implementation of Store, useful for testing.
type MemoryStore struct {
	mu      sync.RWMutex
	secrets map[string]string
}

// NewMemoryStore creates a new empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		secrets: make(map[string]string),
	}
}

// Get retrieves a secret by key.
func (m *MemoryStore) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.secrets[key]
	if !ok {
		return "", ErrSecretNotFound
	}
	return v, nil
}

// Set stores a secret by key.
func (m *MemoryStore) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secrets[key] = value
	return nil
}

// Delete removes a secret by key.
func (m *MemoryStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.secrets[key]; !ok {
		return ErrSecretNotFound
	}
	delete(m.secrets, key)
	return nil
}

// List returns all secret keys in the store.
func (m *MemoryStore) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.secrets))
	for k := range m.secrets {
		keys = append(keys, k)
	}
	return keys, nil
}
