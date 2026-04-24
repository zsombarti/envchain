// Package prompt provides utilities for securely reading sensitive input
// from the terminal, such as passphrases and confirmation prompts.
package prompt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrAborted is returned when the user cancels a prompt.
var ErrAborted = errors.New("prompt: aborted by user")

// ErrMismatch is returned when confirmation input does not match.
var ErrMismatch = errors.New("prompt: inputs do not match")

// Options configures prompt behaviour.
type Options struct {
	// In is the input reader; defaults to os.Stdin.
	In  *os.File
	// Out is where prompts are written; defaults to os.Stderr.
	Out io.Writer
}

func defaults(o *Options) *Options {
	if o == nil {
		o = &Options{}
	}
	if o.In == nil {
		o.In = os.Stdin
	}
	if o.Out == nil {
		o.Out = os.Stderr
	}
	return o
}

// Passphrase prompts the user for a passphrase without echoing input.
// If confirm is true, the user is asked to enter it twice.
func Passphrase(label string, confirm bool, o *Options) (string, error) {
	o = defaults(o)

	fmt.Fprintf(o.Out, "%s: ", label)
	first, err := readPassword(o.In)
	if err != nil {
		return "", err
	}
	fmt.Fprintln(o.Out)

	if !confirm {
		return first, nil
	}

	fmt.Fprintf(o.Out, "Confirm %s: ", label)
	second, err := readPassword(o.In)
	if err != nil {
		return "", err
	}
	fmt.Fprintln(o.Out)

	if first != second {
		return "", ErrMismatch
	}
	return first, nil
}

// Confirm asks a yes/no question and returns true if the user answers "y" or "yes".
func Confirm(question string, o *Options) (bool, error) {
	o = defaults(o)
	fmt.Fprintf(o.Out, "%s [y/N]: ", question)

	var answer string
	if _, err := fmt.Fscanln(o.In, &answer); err != nil && err != io.EOF {
		return false, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func readPassword(f *os.File) (string, error) {
	fd := int(f.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	// Fallback for non-terminal (e.g. piped input in tests).
	var line string
	if _, err := fmt.Fscanln(f, &line); err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}
