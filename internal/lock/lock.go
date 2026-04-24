// Package lock provides session-based locking for envchain profiles.
// A lock prevents concurrent modifications to a profile by recording a
// lock file with a timestamp and optional owner hint.
package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrLocked is returned when a profile is already locked.
var ErrLocked = errors.New("lock: profile is locked by another session")

// ErrNotLocked is returned when attempting to release a lock that does not exist.
var ErrNotLocked = errors.New("lock: profile is not locked")

// Info holds metadata stored inside a lock file.
type Info struct {
	Owner     string    `json:"owner"`
	AcquiredAt time.Time `json:"acquired_at"`
}

// Manager manages lock files for profiles inside a base directory.
type Manager struct {
	dir string
}

// NewManager returns a Manager that stores lock files under dir.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) lockPath(profile string) string {
	return filepath.Join(m.dir, profile+".lock")
}

// Acquire creates a lock for profile. Returns ErrLocked if one already exists.
func (m *Manager) Acquire(profile, owner string) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return fmt.Errorf("lock: mkdir: %w", err)
	}
	path := m.lockPath(profile)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return ErrLocked
		}
		return fmt.Errorf("lock: create: %w", err)
	}
	defer f.Close()
	info := Info{Owner: owner, AcquiredAt: time.Now().UTC()}
	return json.NewEncoder(f).Encode(info)
}

// Release removes the lock for profile. Returns ErrNotLocked if absent.
func (m *Manager) Release(profile string) error {
	err := os.Remove(m.lockPath(profile))
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotLocked
		}
		return fmt.Errorf("lock: remove: %w", err)
	}
	return nil
}

// Status returns the Info for a locked profile, or ErrNotLocked if free.
func (m *Manager) Status(profile string) (*Info, error) {
	f, err := os.Open(m.lockPath(profile))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotLocked
		}
		return nil, fmt.Errorf("lock: open: %w", err)
	}
	defer f.Close()
	var info Info
	if err := json.NewDecoder(f).Decode(&info); err != nil {
		return nil, fmt.Errorf("lock: decode: %w", err)
	}
	return &info, nil
}

// IsLocked reports whether profile is currently locked.
func (m *Manager) IsLocked(profile string) bool {
	_, err := m.Status(profile)
	return err == nil
}
