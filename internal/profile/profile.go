// Package profile manages named environment profiles that group
// related layers and their associated store keys for a project.
package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// ErrNotFound is returned when a profile does not exist.
var ErrNotFound = errors.New("profile not found")

// ErrAlreadyExists is returned when a profile already exists.
var ErrAlreadyExists = errors.New("profile already exists")

// Profile represents a named collection of layer references.
type Profile struct {
	Name   string   `json:"name"`
	Layers []string `json:"layers"`
}

// Manager handles persistence and retrieval of profiles.
type Manager struct {
	dir string
}

// NewManager creates a Manager that stores profiles under dir.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) path(name string) string {
	return filepath.Join(m.dir, name+".json")
}

// Save persists a profile to disk, creating it if it does not exist.
func (m *Manager) Save(p Profile) error {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path(p.Name), data, 0o600)
}

// Load retrieves a profile by name.
func (m *Manager) Load(name string) (Profile, error) {
	data, err := os.ReadFile(m.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return Profile{}, ErrNotFound
	}
	if err != nil {
		return Profile{}, err
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, err
	}
	return p, nil
}

// Delete removes a profile by name.
func (m *Manager) Delete(name string) error {
	err := os.Remove(m.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

// List returns the names of all stored profiles.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
