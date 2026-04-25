// Package watch implements lightweight file-system polling for envchain.
//
// # Overview
//
// The [Watcher] type monitors a list of file paths by comparing modification
// times on a configurable interval. When a tracked file's mtime advances, a
// [Event] is sent on the Watcher.Changes channel so callers can react —
// for example by reloading a profile or re-resolving the active chain.
//
// # Usage
//
//	w := watch.New([]string{"/home/user/.envchain/default.json"}, 0)
//	go w.Run(ctx)
//	for ev := range w.Changes {
//		fmt.Println("changed:", ev.Path)
//	}
//
// # Notes
//
// Polling is used intentionally to keep the implementation simple and
// dependency-free. For most local-dev workflows a 2-second poll interval
// provides adequate responsiveness without measurable overhead.
package watch
