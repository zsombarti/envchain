// Package audit provides a simple event log for tracking environment
// variable access and mutations during a dev session.
package audit

import (
	"fmt"
	"io"
	"os"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventGet    EventKind = "GET"
	EventSet    EventKind = "SET"
	EventDelete EventKind = "DELETE"
	EventResolve EventKind = "RESOLVE"
)

// Event represents a single auditable action.
type Event struct {
	Time    time.Time
	Kind    EventKind
	Key     string
	Source  string
}

// Logger records audit events to a writer.
type Logger struct {
	w      io.Writer
	nowFn  func() time.Time
}

// New creates a Logger that writes to w.
// Pass nil to discard all output.
func New(w io.Writer) *Logger {
	if w == nil {
		w = io.Discard
	}
	return &Logger{w: w, nowFn: time.Now}
}

// NewStderr creates a Logger that writes to stderr.
func NewStderr() *Logger {
	return New(os.Stderr)
}

// Log records an event.
func (l *Logger) Log(kind EventKind, key, source string) Event {
	e := Event{
		Time:   l.nowFn(),
		Kind:   kind,
		Key:    key,
		Source: source,
	}
	fmt.Fprintf(l.w, "%s [%s] key=%q source=%q\n",
		e.Time.Format(time.RFC3339), e.Kind, e.Key, e.Source)
	return e
}

// Get logs a GET event and returns the event.
func (l *Logger) Get(key, source string) Event {
	return l.Log(EventGet, key, source)
}

// Set logs a SET event.
func (l *Logger) Set(key, source string) Event {
	return l.Log(EventSet, key, source)
}

// Delete logs a DELETE event.
func (l *Logger) Delete(key, source string) Event {
	return l.Log(EventDelete, key, source)
}

// Resolve logs a RESOLVE event.
func (l *Logger) Resolve(key, source string) Event {
	return l.Log(EventResolve, key, source)
}
