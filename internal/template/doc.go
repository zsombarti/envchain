// Package template provides lightweight variable interpolation for environment
// value strings used throughout envchain.
//
// Syntax
//
// Two reference forms are supported:
//
//	$VAR        – bare identifier, terminated by a non-identifier character
//	${VAR}      – braced form, allowing adjacent text without ambiguity
//
// A literal dollar sign can be produced with $$.
//
// Expansion is applied iteratively up to Options.MaxDepth times, so values
// that themselves contain references (e.g. BIN=${BASE}/bin where BASE is also
// defined) are resolved transitively.
//
// Usage
//
//	env := map[string]string{"HOST": "localhost", "DSN": "postgres://${HOST}/db"}
//	out, err := template.ExpandMap(env, template.DefaultOptions())
//	// out["DSN"] == "postgres://localhost/db"
package template
