// Package diff computes the difference between two environment variable maps,
// useful for showing what changed between chain resolutions or profile loads.
package diff

import "sort"

// ChangeKind describes the type of change for a single key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level diff entry.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds all changes between two env maps.
type Result struct {
	Changes []Change
}

// Added returns only the added changes.
func (r Result) Added() []Change { return r.filter(Added) }

// Removed returns only the removed changes.
func (r Result) Removed() []Change { return r.filter(Removed) }

// Modified returns only the modified changes.
func (r Result) Modified() []Change { return r.filter(Modified) }

func (r Result) filter(kind ChangeKind) []Change {
	out := make([]Change, 0)
	for _, c := range r.Changes {
		if c.Kind == kind {
			out = append(out, c)
		}
	}
	return out
}

// HasChanges reports whether any non-unchanged entries exist.
func (r Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Compute returns the diff between old and new env maps.
// Keys present in neither are ignored; all keys from both maps are considered.
func Compute(old, next map[string]string) Result {
	keys := unionKeys(old, next)
	changes := make([]Change, 0, len(keys))
	for _, k := range keys {
		ov, inOld := old[k]
		nv, inNew := next[k]
		switch {
		case inOld && !inNew:
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: ov})
		case !inOld && inNew:
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: nv})
		case ov != nv:
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: ov, NewValue: nv})
		default:
			changes = append(changes, Change{Key: k, Kind: Unchanged, OldValue: ov, NewValue: nv})
		}
	}
	return Result{Changes: changes}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
