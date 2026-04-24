// Package run provides functionality to execute a subprocess with a
// resolved environment variable set injected into its environment.
package run

import (
	"errors"
	"os"
	"os/exec"
)

// ErrNoCommand is returned when no command is provided.
var ErrNoCommand = errors.New("run: no command specified")

// Runner executes a command with a given environment overlay.
type Runner struct {
	// Inherit controls whether the host environment is inherited.
	Inherit bool
}

// New returns a Runner that inherits the host environment by default.
func New() *Runner {
	return &Runner{Inherit: true}
}

// Exec runs the given command with args, injecting env as additional
// environment variables. If Inherit is true, the host environment is
// merged first, with env values taking precedence.
//
// The function blocks until the process exits and returns any error
// from the subprocess (including non-zero exit codes via *exec.ExitError).
func (r *Runner) Exec(env map[string]string, command string, args ...string) error {
	if command == "" {
		return ErrNoCommand
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if r.Inherit {
		cmd.Env = os.Environ()
	}

	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	return cmd.Run()
}

// EnvSlice converts a map of key/value pairs to a slice of KEY=VALUE strings.
func EnvSlice(env map[string]string) []string {
	slice := make([]string, 0, len(env))
	for k, v := range env {
		slice = append(slice, k+"="+v)
	}
	return slice
}
