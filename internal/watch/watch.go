// Package watch provides file-system watching utilities that notify
// callers when tracked environment files are modified on disk.
package watch

import (
	"context"
	"os"
	"sync"
	"time"
)

// Event is emitted when a watched file changes.
type Event struct {
	Path    string
	ModTime time.Time
}

// Watcher polls a set of file paths and sends an Event on Changes whenever
// a file's modification time advances.
type Watcher struct {
	paths    []string
	interval time.Duration
	Changes  chan Event

	mu      sync.Mutex
	mtimes  map[string]time.Time
}

// New creates a Watcher that polls the given paths at the given interval.
// A zero interval defaults to 2 seconds.
func New(paths []string, interval time.Duration) *Watcher {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	w := &Watcher{
		paths:    paths,
		interval: interval,
		Changes:  make(chan Event, len(paths)+1),
		mtimes:   make(map[string]time.Time),
	}
	// seed initial mtimes so first poll does not fire spuriously
	for _, p := range paths {
		if fi, err := os.Stat(p); err == nil {
			w.mtimes[p] = fi.ModTime()
		}
	}
	return w
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll()
		}
	}
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, p := range w.paths {
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		prev, seen := w.mtimes[p]
		if !seen || fi.ModTime().After(prev) {
			w.mtimes[p] = fi.ModTime()
			if seen {
				select {
				case w.Changes <- Event{Path: p, ModTime: fi.ModTime()}:
				default:
				}
			}
		}
	}
}
