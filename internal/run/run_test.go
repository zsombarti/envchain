package run_test

import (
	"os/exec"
	"sort"
	"testing"

	"envchain/internal/run"
)

func TestExecNoCommand(t *testing.T) {
	r := run.New()
	err := r.Exec(nil, "")
	if err != run.ErrNoCommand {
		t.Fatalf("expected ErrNoCommand, got %v", err)
	}
}

func TestExecSuccess(t *testing.T) {
	r := run.New()
	err := r.Exec(nil, "true")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestExecNonZeroExit(t *testing.T) {
	r := run.New()
	err := r.Exec(nil, "false")
	if err == nil {
		t.Fatal("expected non-nil error for exit code 1")
	}
	var exitErr *exec.ExitError
	if ok := isExitError(err, &exitErr); !ok {
		t.Fatalf("expected *exec.ExitError, got %T", err)
	}
}

func TestExecInjectsEnv(t *testing.T) {
	r := run.New()
	env := map[string]string{"ENVCHAIN_TEST_VAR": "hello"}
	// sh -c 'test "$ENVCHAIN_TEST_VAR" = hello' exits 0 on match
	err := r.Exec(env, "sh", "-c", `test "$ENVCHAIN_TEST_VAR" = hello`)
	if err != nil {
		t.Fatalf("env var not injected: %v", err)
	}
}

func TestExecNoInherit(t *testing.T) {
	r := &run.Runner{Inherit: false}
	// Without inheriting PATH we still use absolute path.
	err := r.Exec(nil, "/usr/bin/true")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestEnvSlice(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	slice := run.EnvSlice(env)
	sort.Strings(slice)
	if len(slice) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(slice))
	}
	if slice[0] != "A=1" || slice[1] != "B=2" {
		t.Fatalf("unexpected slice: %v", slice)
	}
}

// isExitError is a helper that avoids importing errors in test file.
func isExitError(err error, target **exec.ExitError) bool {
	e, ok := err.(*exec.ExitError)
	if ok {
		*target = e
	}
	return ok
}
