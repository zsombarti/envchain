// Package resolve provides higher-level environment variable resolution
// built on top of a chain.Chain.
//
// It supports:
//
//   - Prefix filtering — only keys that start with a given prefix are included
//     in the result.
//   - Prefix stripping — the common prefix can be removed from key names so
//     consumers receive clean variable names.
//   - Overrides — a caller-supplied map is merged last, giving it the highest
//     precedence over any layer in the chain.
//
// Typical usage:
//
//	env, err := resolve.Resolve(myChain, resolve.Options{
//		Prefix:      "APP_",
//		StripPrefix: true,
//		Overrides:   map[string]string{"DEBUG": "true"},
//	})
package resolve
