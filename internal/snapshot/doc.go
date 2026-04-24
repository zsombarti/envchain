// Package snapshot provides functionality to capture, persist, compare,
// and manage named snapshots of resolved environment variable sets.
//
// A Snapshot records the full set of key-value pairs at a specific moment,
// along with a label and creation timestamp. Snapshots can be saved to and
// loaded from JSON files on disk.
//
// The Manager type provides a higher-level interface for storing multiple
// named snapshots in a directory, listing available snapshots, and deleting
// them by label.
//
// The Diff function compares two snapshots and reports which keys were added,
// removed, or changed between them — useful for auditing environment drift
// across different workflow stages.
//
// Example usage:
//
//	mgr, _ := snapshot.NewManager(".envchain/snapshots")
//	snap := snapshot.New("pre-deploy", resolvedVars)
//	mgr.Save(snap)
//
//	prev, _ := mgr.Load("pre-deploy")
//	added, removed, changed := snapshot.Diff(prev, snap)
package snapshot
