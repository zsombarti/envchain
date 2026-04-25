// Package diff provides utilities for computing and displaying the difference
// between two environment variable maps.
//
// It is used by envchain to show what changes occur when switching profiles,
// applying a chain layer, or comparing snapshots.
//
// Basic usage:
//
//	old := map[string]string{"FOO": "bar"}
//	next := map[string]string{"FOO": "baz", "NEW": "val"}
//
//	r := diff.Compute(old, next)
//	diff.Write(os.Stdout, r, diff.DefaultOptions())
//	fmt.Println(diff.Summary(r))
//
// Output:
//
//	~ FOO: bar -> baz
//	+ NEW=val
//	~1 +1
package diff
