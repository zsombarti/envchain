// Package cache provides a time-bounded in-memory cache for resolved
// environment variable sets, reducing repeated secret-store lookups during
// a single dev session.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached map of environment variables and the time it expires.
type Entry struct {
	Vars      map[string]string
	ExpiresAt time.Time
}

// Cache is a thread-safe, TTL-based store for resolved env maps keyed by
// profile name.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
	now     func() time.Time
}

// New creates a Cache with the given TTL. Entries older than ttl are
// considered stale and will not be returned.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
		now:     time.Now,
	}
}

// Set stores a copy of vars under key, expiring after the configured TTL.
func (c *Cache) Set(key string, vars map[string]string) {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = Entry{
		Vars:      copy,
		ExpiresAt: c.now().Add(c.ttl),
	}
}

// Get returns the cached vars for key and true if a valid, non-expired entry
// exists. Otherwise it returns nil and false.
func (c *Cache) Get(key string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || c.now().After(e.ExpiresAt) {
		return nil, false
	}
	copy := make(map[string]string, len(e.Vars))
	for k, v := range e.Vars {
		copy[k] = v
	}
	return copy, true
}

// Delete removes the entry for key, if present.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]Entry)
}

// Len returns the number of entries currently held (including stale ones).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
