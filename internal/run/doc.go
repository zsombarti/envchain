// Package run provides subprocess execution with environment injection
// for the envchain CLI tool.
//
// # Overview
//
// The Runner type wraps [os/exec.Cmd] and allows callers to supply a
// map of environment variables that are merged into (or replace) the
// host process environment before the child process is started.
//
// # Usage
//
//	r := run.New()          // inherits host env by default
//	env := map[string]string{"DB_URL": "postgres://..."}
//	err := r.Exec(env, "go", "test", "./...")
//
// Setting Runner.Inherit to false gives the subprocess a clean
// environment containing only the variables supplied by the caller.
package run
