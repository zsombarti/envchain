package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"envchain/internal/watch"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestWatcherDetectsChange(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "env.json")
	writeFile(t, p, `{"A":"1"}`)

	w := watch.New([]string{p}, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go w.Run(ctx)

	// give the watcher one tick to seed mtimes
	time.Sleep(30 * time.Millisecond)

	// modify the file
	writeFile(t, p, `{"A":"2"}`)

	select {
	case ev := <-w.Changes:
		if ev.Path != p {
			t.Fatalf("expected path %q, got %q", p, ev.Path)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcherNoSpuriousEventOnStart(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "env.json")
	writeFile(t, p, `{}`)

	w := watch.New([]string{p}, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go w.Run(ctx)

	// no writes — Changes should stay empty
	<-ctx.Done()
	select {
	case ev := <-w.Changes:
		t.Fatalf("unexpected event: %+v", ev)
	default:
	}
}

func TestWatcherMissingFileIgnored(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.json")

	w := watch.New([]string{missing}, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go w.Run(ctx)
	<-ctx.Done()

	select {
	case ev := <-w.Changes:
		t.Fatalf("unexpected event for missing file: %+v", ev)
	default:
	}
}

func TestWatcherMultiplePaths(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.json")
	p2 := filepath.Join(dir, "b.json")
	writeFile(t, p1, `{}`)
	writeFile(t, p2, `{}`)

	w := watch.New([]string{p1, p2}, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go w.Run(ctx)
	time.Sleep(30 * time.Millisecond)

	writeFile(t, p2, `{"X":"1"}`)

	select {
	case ev := <-w.Changes:
		if ev.Path != p2 {
			t.Fatalf("expected %q, got %q", p2, ev.Path)
		}
	case <-ctx.Done():
		t.Fatal("timed out")
	}
}
