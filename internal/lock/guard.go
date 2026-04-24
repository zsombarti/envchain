package lock

import "fmt"

// Guard holds an acquired lock and releases it when Done is called.
// Obtain a Guard via [WithLock].
type Guard struct {
	m       *Manager
	profile string
}

// Done releases the lock. It is safe to call Done multiple times;
// subsequent calls return ErrNotLocked which is silently ignored.
func (g *Guard) Done() error {
	err := g.m.Release(g.profile)
	if err == ErrNotLocked {
		return nil
	}
	return err
}

// WithLock acquires a lock for profile and returns a Guard that releases it.
// Returns ErrLocked if the profile is already locked.
//
//	guard, err := lock.WithLock(m, "staging", "deploy-bot")
//	if err != nil { ... }
//	defer guard.Done()
func WithLock(m *Manager, profile, owner string) (*Guard, error) {
	if err := m.Acquire(profile, owner); err != nil {
		return nil, fmt.Errorf("lock: withlock %q: %w", profile, err)
	}
	return &Guard{m: m, profile: profile}, nil
}
