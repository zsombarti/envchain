// Package profile provides management of named environment profiles for
// envchain. A profile is a lightweight descriptor that records an ordered
// list of layer names to be resolved by the chain engine.
//
// Profiles are stored as JSON files in a configurable directory, one file
// per profile, making them easy to inspect, diff, and version-control.
//
// Typical usage:
//
//	m := profile.NewManager("~/.config/envchain/profiles")
//
//	// Create a new profile
//	_ = m.Save(profile.Profile{
//		Name:   "dev",
//		Layers: []string{"global", "project", "local"},
//	})
//
//	// Load and inspect
//	p, _ := m.Load("dev")
//	fmt.Println(p.Layers)
//
//	// List all profiles
//	names, _ := m.List()
package profile
