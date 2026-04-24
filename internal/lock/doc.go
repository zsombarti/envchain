// Package lock implements advisory file-based locking for envchain profiles.
//
// A lock file is created under a configurable directory when a profile is
// acquired. The file stores a JSON-encoded [Info] record containing the owner
// hint and the acquisition timestamp.
//
// Typical usage:
//
//	m := lock.NewManager("/home/user/.config/envchain/locks")
//
//	if err := m.Acquire("production", os.Getenv("USER")); err != nil {
//		if errors.Is(err, lock.ErrLocked) {
//			info, _ := m.Status("production")
//			log.Fatalf("profile locked by %s since %s", info.Owner, info.AcquiredAt)
//		}
//		log.Fatal(err)
//	}
//	defer m.Release("production")
//
// Locks are cooperative — they rely on callers honouring [ErrLocked] rather
// than OS-level exclusive locks. This is intentional: it keeps the
// implementation portable and avoids stale-lock issues across NFS mounts.
package lock
