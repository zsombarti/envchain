package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the resolved environment variables at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Vars      map[string]string `json:"vars"`
}

// New creates a new Snapshot with the given label and resolved vars.
func New(label string, vars map[string]string) *Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Vars:      copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %q: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Diff returns keys that differ between two snapshots.
// Returns added, removed, and changed key sets.
func Diff(a, b *Snapshot) (added, removed, changed []string) {
	for k, bv := range b.Vars {
		if av, ok := a.Vars[k]; !ok {
			added = append(added, k)
		} else if av != bv {
			changed = append(changed, k)
		}
	}
	for k := range a.Vars {
		if _, ok := b.Vars[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}
