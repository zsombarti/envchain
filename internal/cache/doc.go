// Package cache implements a lightweight, TTL-based in-memory cache for
// resolved environment variable maps.
//
// During a local dev session the same profile may be resolved many times
// (e.g. on every file-watch reload). Hitting the secret store on every
// resolution adds latency and may trigger rate-limits. Cache sits between
// the resolver and the store, returning a fresh copy of the last resolved
// map until the entry expires.
//
// # Usage
//
//	c := cache.New(2 * time.Minute)
//
//	// populate after a real resolve
//	c.Set("myprofile", resolvedVars)
//
//	// subsequent calls within TTL skip the store
//	if vars, ok := c.Get("myprofile"); ok {
//		// use vars
//	}
//
// All methods are safe for concurrent use.
package cache
