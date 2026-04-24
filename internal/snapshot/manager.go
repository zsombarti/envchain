package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manager handles storing and retrieving named snapshots from a directory.
type Manager struct {
	dir string
}

// NewManager creates a Manager that persists snapshots under dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot manager: mkdir %q: %w", dir, err)
	}
	return &Manager{dir: dir}, nil
}

// Save persists the snapshot under its label name.
func (m *Manager) Save(s *Snapshot) error {
	path := m.pathFor(s.Label)
	return s.Save(path)
}

// Load retrieves a snapshot by label.
func (m *Manager) Load(label string) (*Snapshot, error) {
	return Load(m.pathFor(label))
}

// Delete removes a snapshot by label.
func (m *Manager) Delete(label string) error {
	path := m.pathFor(label)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot manager: %q not found", label)
		}
		return fmt.Errorf("snapshot manager: delete %q: %w", label, err)
	}
	return nil
}

// List returns all snapshot labels stored in the manager directory.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot manager: list: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			labels = append(labels, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return labels, nil
}

func (m *Manager) pathFor(label string) string {
	return filepath.Join(m.dir, label+".json")
}
